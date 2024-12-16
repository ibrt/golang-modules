package tpgm

import (
	"context"
	"embed"
	"fmt"
	"net/url"
	"strings"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/idz"
	"github.com/ibrt/golang-utils/memz"
	"github.com/ibrt/golang-utils/urlz"
	"github.com/jackc/pgx/v5"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/pgm"
)

var (
	_ fixturez.BeforeSuite = (*Helper)(nil)
	_ fixturez.AfterSuite  = (*Helper)(nil)
	_ fixturez.BeforeTest  = (*Helper)(nil)
)

// Helper is a test helper.
type Helper struct {
	MigrationsFS    *embed.FS
	dbName          string
	origPostgresURL string
	releaser        func()
}

// BeforeSuite implements [fixturez.BeforeSuite].
func (h *Helper) BeforeSuite(ctx context.Context, _ *gomega.WithT) context.Context {
	cfg := cfgm.MustGet[pgm.PGConfigMixin](ctx).GetPGConfig()
	h.dbName = fmt.Sprintf("%v-db", idz.MustNewRandomUUID())

	h.mustCreateDB(ctx, cfg.PostgresURL, h.dbName)
	h.origPostgresURL, cfg.PostgresURL = cfg.PostgresURL, h.mustSelectDB(cfg.PostgresURL, h.dbName)

	injector, releaser := pgm.NewInitializer(h.MigrationsFS)(ctx)
	h.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements [fixturez.AfterSuite].
func (h *Helper) AfterSuite(ctx context.Context, _ *gomega.WithT) {
	h.releaser()
	h.releaser = nil

	cfg := cfgm.MustGet[pgm.PGConfigMixin](ctx).GetPGConfig()
	cfg.PostgresURL = h.origPostgresURL
	h.mustDropDB(ctx, cfg.PostgresURL, h.dbName)

	h.origPostgresURL = ""
	h.dbName = ""
}

// BeforeTest implements [fixturez.BeforeTest].
func (h *Helper) BeforeTest(ctx context.Context, _ *gomega.WithT, _ *gomock.Controller) context.Context {
	type table struct {
		TableName string `db:"tablename"`
	}

	rows, err := pgm.MustGet(ctx).Query("test.list_tables",
		`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename <> $1`,
		pgm.MigrationsTableName)
	errorz.MaybeMustWrap(err)

	tables, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[table])
	errorz.MaybeMustWrap(err)

	_, err = pgm.MustGet(ctx).Exec("test.truncate",
		fmt.Sprintf(
			`TRUNCATE TABLE %v RESTART IDENTITY RESTRICT`,
			strings.Join(memz.TransformSlice(tables, func(_ int, t *table) string { return pgx.Identifier{t.TableName}.Sanitize() }), ", ")))
	errorz.MaybeMustWrap(err)

	return ctx
}

func (h *Helper) mustOpen(ctx context.Context, pgURL string) *pgx.Conn {
	pg, err := pgx.Connect(ctx, pgURL)
	errorz.MaybeMustWrap(err)
	return pg
}

func (h *Helper) mustCreateDB(ctx context.Context, pgURL, dbName string) {
	dbName = pgx.Identifier{dbName}.Sanitize()
	pg := h.mustOpen(ctx, pgURL)
	defer func() {
		_ = pg.Close(context.Background())
	}()

	_, err := pg.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %v`, dbName))
	errorz.MaybeMustWrap(err)

	_, err = pg.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE %v`, dbName))
	errorz.MaybeMustWrap(err)
}

func (h *Helper) mustDropDB(ctx context.Context, pgURL, dbName string) {
	dbName = pgx.Identifier{dbName}.Sanitize()
	pg := h.mustOpen(ctx, pgURL)
	defer func() {
		_ = pg.Close(context.Background())
	}()

	_, err := pg.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %v`, dbName))
	errorz.MaybeMustWrap(err)
}

func (h *Helper) mustSelectDB(pgURL, dbName string) string {
	return urlz.MustEdit(pgURL, func(u *url.URL) {
		u.Path = dbName
	})
}
