package rules

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"log/slog"

	"emperror.dev/errors"
	glob "github.com/bmatcuk/doublestar/v4"
	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/logger"
)

type FilesystemRuleFinder struct {
	staticRuleFinderOptions
	config *v1beta1.FilesystemRulesSpec
	logger *slog.Logger
}

type staticRuleFinderOptions struct {
	fs fs.FS
}

type FilesystemRuleFinderOption func(*staticRuleFinderOptions)

func (o *staticRuleFinderOptions) apply(opts ...FilesystemRuleFinderOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithFS(fs fs.FS) FilesystemRuleFinderOption {
	return func(o *staticRuleFinderOptions) {
		o.fs = fs
	}
}

func NewFilesystemRuleFinder(config *v1beta1.FilesystemRulesSpec, opts ...FilesystemRuleFinderOption) *FilesystemRuleFinder {
	options := staticRuleFinderOptions{
		fs: os.DirFS("/"),
	}
	options.apply(opts...)

	return &FilesystemRuleFinder{
		staticRuleFinderOptions: options,
		config:                  config,
		logger:                  logger.New().WithGroup("rules"),
	}
}

func (f *FilesystemRuleFinder) Find(context.Context) ([]RuleGroup, error) {
	groups := []rulefmt.RuleGroup{}

	for _, pathExpr := range f.config.PathExpressions {
		pathExpr = strings.TrimPrefix(pathExpr, "/")
		matched, err := glob.Glob(f.fs, pathExpr)
		lg := f.logger.With("expression", pathExpr)
		if err != nil {
			lg.Warn("error searching for rules files using path expression", logger.Err(err))

			continue
		}

		lg.Debug(fmt.Sprintf("found %d rules files matching path expression", len(matched)))
		for _, path := range matched {
			lg := lg.With("path", path)
			data, err := fs.ReadFile(f.fs, path)
			if err != nil {
				lg.Warn("error reading rules file", logger.Err(err))

				continue
			}
			list, errs := rulefmt.Parse(data)
			if len(errs) > 0 {
				lg.Warn("error parsing rules file", logger.Err(errors.Combine(errs...)))

				continue
			}
			groups = append(groups, list.Groups...)
			f.logger.Debug(fmt.Sprintf("found %d rule groups in file %s", len(list.Groups), path))
		}
	}

	f.logger.Info(fmt.Sprintf("found %d rule groups in filesystem", len(groups)))
	ruleGroups := []RuleGroup{}
	for _, g := range groups {
		ruleGroups = append(ruleGroups, RuleGroup(g))
	}
	return ruleGroups, nil
}
