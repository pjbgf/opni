package targets

import (
	"context"

	"github.com/magefile/mage/mg"
)

func (Build) Opni(ctx context.Context) error {
	mg.CtxDeps(ctx, Build.Archives)

	_, tr := Tracer.Start(ctx, "target.build.opni")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/opni",
		Output: "bin/opni",
		Tags:   []string{"nomsgpack"},
	})
}

func (Build) OpniMinimal(ctx context.Context) error {
	_, tr := Tracer.Start(ctx, "target.build.opni-minimal")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/opni",
		Output: "bin/opni-minimal",
		Tags:   []string{"nomsgpack", "minimal"},
	})
}
