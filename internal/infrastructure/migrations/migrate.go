package migrations

import (
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/migrations/versions"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// Migration represents a migration with Up and Down functions
// Migration defines the structure for each migration, which includes
// an Up function to apply the migration and a Down function to revert it.
type Migration struct {
	Version string
	Up      func(*mongo.Database) error // Up function to apply the migration
	Down    func(*mongo.Database) error // Down function to revert the migration
}

// List of available migrations
// The Migrations slice stores all the migrations, along with their Up and Down functions
var Migrations = []Migration{
	{"000", versions.Up000, versions.Down000},
	{"001", versions.Up001, versions.Down001},
}

// Finds the index of a migration by version
// The function searches for the migration index in the Migrations slice by its version.
func findMigrationIndex(version string) (int, error) {
	for i, migration := range Migrations {
		if migration.Version == version {
			return i, nil
		}
	}
	return -1, fmt.Errorf("version %s not found", version)
}

// Runs migrations up to a specified target version (Up)
// This function will apply all migrations up to the specified target version.
func runMigrate(db *mongo.Database, targetVersion string) error {
	// Find the target migration index
	targetIndex, err := findMigrationIndex(targetVersion)
	if err != nil {
		return err
	}

	// Loop through migrations and apply them if they haven't been applied
	for i := 0; i <= targetIndex; i++ {
		migration := Migrations[i]
		// Check if the migration has already been applied
		applied, err := HasMigration(db, migration.Version)
		if err != nil {
			return fmt.Errorf("error checking migration status %s: %w", migration.Version, err)
		}

		// Apply the migration if it has not been applied yet
		if !applied {
			log.Printf("Running migration UP: %s", migration.Version)
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("error applying migration %s: %w", migration.Version, err)
			}
			// Log the applied migration
			if err := AddMigrationLog(db, migration.Version); err != nil {
				return fmt.Errorf("error logging migration %s: %w", migration.Version, err)
			}
		} else {
			log.Printf("Migration %s already applied, skipping...", migration.Version)
		}
	}

	log.Printf("Migrations applied up to version %s", targetVersion)
	return nil
}

// Rolls back migrations down to a specified target version (Down)
// This function will roll back migrations starting from the latest one
// until it reaches the target version.
func runUnmigrate(db *mongo.Database, targetVersion string) error {
	// Find the target migration index
	targetIndex, err := findMigrationIndex(targetVersion)
	if err != nil {
		return err
	}

	// Loop through migrations in reverse order and roll them back
	for i := len(Migrations) - 1; i > targetIndex; i-- {
		migration := Migrations[i]
		// Check if the migration has been applied
		applied, err := HasMigration(db, migration.Version)
		if err != nil {
			return fmt.Errorf("error checking migration status %s: %w", migration.Version, err)
		}

		// Roll back the migration if it has been applied
		if applied {
			log.Printf("Running migration DOWN: %s", migration.Version)
			if err := migration.Down(db); err != nil {
				return fmt.Errorf("error rolling back migration %s: %w", migration.Version, err)
			}
			// Remove the migration log
			if err := RemoveMigrationLog(db, migration.Version); err != nil {
				return fmt.Errorf("error removing migration log %s: %w", migration.Version, err)
			}
		} else {
			log.Printf("Migration %s not applied, skipping...", migration.Version)
		}
	}

	log.Printf("Migrations rolled back down to version %s", targetVersion)
	return nil
}

// Creates CLI commands for migrations
// This function defines the available CLI commands for migrations, such as
// migrating, rolling back migrations, and showing applied migrations.
func migrateCommands(client *mongo.Client, config *config.AppConfig) []cli.Command {
	db := client.Database(config.MongoDB)

	// Define flags for specifying the target version
	// The "version" flag specifies the target version for migrations
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Target migration version (e.g., 001, 002)",
			Value:   Migrations[len(Migrations)-1].Version, // Default to the latest version
		},
	}

	// Command to apply migrations
	// The "migrate" command applies migrations up to the specified version
	migrate := &cli.Command{
		Name:  "migrate",
		Usage: "Apply migrations up to a specified version",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			targetVersion := cCtx.String("version")

			if err := runMigrate(db, targetVersion); err != nil {
				log.Fatalf("Error applying migrations: %v", err)
				return err
			}
			log.Println("Migrations applied successfully")
			return nil
		},
	}

	// Command to roll back migrations
	// The "unmigrate" command rolls back migrations down to the specified version
	unmigrate := &cli.Command{
		Name:  "unmigrate",
		Usage: "Rollback migrations down to a specified version",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			targetVersion := cCtx.String("version")

			if err := runUnmigrate(db, targetVersion); err != nil {
				log.Fatalf("Error rolling back migrations: %v", err)
				return err
			}
			log.Println("Migrations rolled back successfully")
			return nil
		},
	}

	// Command to show applied migrations
	// The "showmigrations" command lists all migrations that have been applied to the database
	showMigrations := &cli.Command{
		Name:  "showmigrations",
		Usage: "Show all migrations that have been applied",
		Action: func(cCtx *cli.Context) error {
			// Get the list of applied migrations from the DB
			appliedMigrations := []string{}
			for _, migration := range Migrations {
				applied, err := HasMigration(db, migration.Version)
				if err != nil {
					log.Fatalf("Error checking migration status %s: %v", migration.Version, err)
					return err
				}

				// Add applied migration to the list
				if applied {
					appliedMigrations = append(appliedMigrations, migration.Version)
				}
			}

			// Display the list of applied migrations
			if len(appliedMigrations) > 0 {
				log.Println("Applied migrations:")
				for _, version := range appliedMigrations {
					log.Printf("- %s", version)
				}
			} else {
				log.Println("No migrations have been applied yet.")
			}

			return nil
		},
	}

	// Return all the defined migration-related commands
	return []cli.Command{*migrate, *unmigrate, *showMigrations}
}

type MigrationsCommands struct {
	Commands []cli.Command
}

// Public function to get the migration CLI commands
// This function serves as an entry point for obtaining the migration commands.
func NewMigrateCommands(client *mongo.Client, config *config.AppConfig) MigrationsCommands {
	return MigrationsCommands{Commands: migrateCommands(client, config)}
}
