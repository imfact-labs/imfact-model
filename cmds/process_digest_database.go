package cmds

import (
	"context"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	mongodbst "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	"github.com/ProtoconNet/mitum-minic/digest"
	"github.com/ProtoconNet/mitum2/isaac"
	isaacdatabase "github.com/ProtoconNet/mitum2/isaac/database"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/pkg/errors"
)

func ProcessDatabase(ctx context.Context) (context.Context, error) {
	var l cdigest.YamlDigestDesign
	if err := util.LoadFromContext(ctx, cdigest.ContextValueDigestDesign, &l); err != nil {
		return ctx, err
	}

	if l.Equal(cdigest.YamlDigestDesign{}) {
		return ctx, nil
	}
	conf := l.Database()

	switch {
	case conf.URI().Scheme == "mongodb", conf.URI().Scheme == "mongodb+srv":
		return processMongodbDatabase(ctx, l)
	default:
		return ctx, errors.Errorf("Unsupported database type, %v", conf.URI().Scheme)
	}
}

func processMongodbDatabase(ctx context.Context, l cdigest.YamlDigestDesign) (context.Context, error) {
	conf := l.Database()

	/*
		ca, err := cache.NewCacheFromURI(conf.Cache().String())
		if err != nil {
			return ctx, err
		}
	*/

	var encs *encoder.Encoders
	if err := util.LoadFromContext(ctx, launch.EncodersContextKey, &encs); err != nil {
		return ctx, err
	}

	st, err := mongodbst.NewDatabaseFromURI(conf.URI().String(), encs)
	if err != nil {
		return ctx, err
	}

	if err := st.Initialize(); err != nil {
		return ctx, err
	}

	var db isaac.Database
	if err := util.LoadFromContextOK(ctx, launch.CenterDatabaseContextKey, &db); err != nil {
		return ctx, err
	}

	mst, ok := db.(*isaacdatabase.Center)
	if !ok {
		return ctx, errors.Errorf("expected isaacdatabase.Center, not %T", db)
	}

	dst, err := loadDigestDatabase(mst, st, false)
	if err != nil {
		return ctx, err
	}
	var log *logging.Logging
	if err := util.LoadFromContextOK(ctx, launch.LoggingContextKey, &log); err != nil {
		return ctx, err
	}

	_ = dst.SetLogging(log)

	return context.WithValue(ctx, cdigest.ContextValueDigestDatabase, dst), nil
}

func loadDigestDatabase(mst *isaacdatabase.Center, st *mongodbst.Database, readonly bool) (*cdigest.Database, error) {
	var dst *cdigest.Database
	if readonly {
		s, err := cdigest.NewReadonlyDatabase(mst, st)
		if err != nil {
			return nil, err
		}
		dst = s
	} else {
		s, err := cdigest.NewDatabase(mst, st)
		if err != nil {
			return nil, err
		}
		dst = s
	}

	if err := dst.Initialize(digest.DefaultIndexes); err != nil {
		return nil, err
	}

	return dst, nil
}
