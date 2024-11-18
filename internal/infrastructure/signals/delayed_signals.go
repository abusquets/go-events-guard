package signals

import (
	"container/list"
	"eventsguard/internal/infrastructure/mylog"
	"fmt"
	"regexp"
	"sync"
)

// Tipus Callback defineix la signatura de les funcions manejadores d'esdeveniments.
type Callback func(args []interface{}) error

type delayedSignalManager struct {
	subscriptions map[string]*list.List // Subscripcions per cada tema.
	eventQueue    *list.List            // Cua d'esdeveniments.
	mu            sync.Mutex            // Mutex per sincronitzar operacions.
	logger        mylog.Logger
	eventChannel  chan struct{} // Canal per notificar els esdeveniments
}

// NewSignalsBus crea una nova instància de delayedSignalManager amb canal.
func NewSignalsBus() SignalsBus {
	return &delayedSignalManager{
		subscriptions: make(map[string]*list.List),
		eventQueue:    list.New(),
		logger:        mylog.GetLogger(),
		eventChannel:  make(chan struct{}, 1), // Canal amb una capacitat de 1 per evitar el bloqueig
	}
}

func (ds *delayedSignalManager) EventChannel() chan struct{} {
	return ds.eventChannel
}

// Subscribe afegeix un callback per un tema específic.
func (ds *delayedSignalManager) Subscribe(topic string, callback Callback) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.subscriptions[topic] == nil {
		ds.subscriptions[topic] = list.New()
	}
	ds.subscriptions[topic].PushBack(callback)
	ds.logger.Info(fmt.Sprintf("Callback registrat per al tema: %s", topic))
}

// Emit afegeix un esdeveniment a la cua utilitzant patrons amb wildcards.
func (ds *delayedSignalManager) AfterTransaction(topic string, args ...interface{}) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Comprovem si el tema té algun patró de subscripció que hi coincideixi
	if !ds.hasMatchingSubscription(topic) {
		return fmt.Errorf("no hi ha subscripcions per al tema: %s", topic)
	}

	// Afegeix l'esdeveniment a la cua
	event := map[string]interface{}{
		"topic": topic,
		"args":  args,
	}
	ds.eventQueue.PushBack(event)
	ds.logger.Info(fmt.Sprintf("Esdeveniment afegit a la cua per al tema: %s", topic))
	return nil
}

// Emit executa immediatament els callbacks per un tema específic.
func (ds *delayedSignalManager) Emit(topic string, args ...interface{}) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	return ds.executeCallbacks(topic, args)
}

// ProcessQueue processa tots els esdeveniments de la cua.
func (ds *delayedSignalManager) ProcessQueue() error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	fmt.Println("###################################ProcessQueue")

	for ds.eventQueue.Len() > 0 {
		// Extraiem el primer esdeveniment de la cua
		elem := ds.eventQueue.Front()
		if elem == nil {
			break
		}

		// Elimina l'esdeveniment de la cua
		ds.eventQueue.Remove(elem)

		// Obtenim el tema i els arguments de l'esdeveniment
		event := elem.Value.(map[string]interface{})
		topic := event["topic"].(string)
		args := event["args"].([]interface{})

		// Executem els callbacks associats
		if err := ds.executeCallbacks(topic, args); err != nil {
			ds.logger.Error(fmt.Sprintf("Error processant esdeveniment per al tema %s: %v", topic, err))
		}
	}

	return nil
}

// executeCallbacks executa els callbacks associats amb un tema donat, utilitzant wildcards.
func (ds *delayedSignalManager) executeCallbacks(topic string, args []interface{}) error {
	fmt.Println("###################################executeCallbacks")
	// Itera per cada patró de subscripció per comprovar si coincideix amb el tema
	for pattern, callbacks := range ds.subscriptions {
		if ds.MatchTopic(pattern, topic) {
			for cb := callbacks.Front(); cb != nil; cb = cb.Next() {
				callback := cb.Value.(Callback)
				if err := callback(args); err != nil {
					ds.logger.Error(fmt.Sprintf("Error processant callback per al tema %s: %v", topic, err))
					return err
				}
				ds.logger.Info(fmt.Sprintf("Callback processat correctament per al tema %s", topic))
			}
		}
	}
	return nil
}

// hasMatchingSubscription verifica si hi ha alguna subscripció que coincideixi amb el tema utilitzant wildcards.
func (ds *delayedSignalManager) hasMatchingSubscription(topic string) bool {
	for pattern := range ds.subscriptions {
		if ds.MatchTopic(pattern, topic) {
			return true
		}
	}
	return false
}

// MatchTopic verifica si un tema coincideix amb un patró (suporta wildcards).
func (ds *delayedSignalManager) MatchTopic(pattern, topic string) bool {
	// Converteix el patró de wildcard en una expressió regular
	pattern = "^" + regexp.QuoteMeta(pattern) + "$"
	pattern = regexp.MustCompile(`\\\*`).ReplaceAllString(pattern, `[^/]+`)
	pattern = regexp.MustCompile(`\\\#`).ReplaceAllString(pattern, `.*`)

	match, _ := regexp.MatchString(pattern, topic)
	return match
}
