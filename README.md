## TypeScript types generator for pocketbase

### Usage

In your pocketbase project, run:

```
go get github.com/davenh99/pb-typescript
```

Then, in your main.go, put in something like (you can use an env var to avoid attaching the plugin in production):

```
if env == "development" {
    gentypes.Register(app, gentypes.Config{
        FilePath:                   "ui",
        PrintSelectOptions:         true,
    })
}
```

The package will look for the root of your git repo (git is required for this package), and then go into the specified file path.

The following files will be generated:

```
base.d.ts
pocketbase-types.ts
select-options.ts // optional
```

Your types will then be generated in the specified directory, whenever the schema changes or you run the following command (assuming your app executable is named pocketbase):

```
./pocketbase gen-types
```

For more examples of how the types look when they are generated check [progressa](https://github.com/davenh99/progressa/blob/main/ui/base.d.ts)

Here is a sample output (the preview is a computed field, attached using [pb-computedfields](https://github.com/davenh99/pb-computedfields):

```
interface BaseRecord {
  readonly id: string;
  readonly collectionName: string;
  readonly collectionId: string;
  readonly created: string;
  readonly updated: string;
}

/* Collection type: base */
interface Routines {
  name: string; // text
  user: string; // relation
  description: string; // text
  exercisesOrder: any; // json
  readonly preview: string; // text
}
type RoutinesRecord = Routines & BaseRecord;
type RoutinesUpdatePayload = Partial<RoutinesRecord>;
```

Why this, and not:
- [pocketbase-types-generator](https://github.com/wiezmankimchi/pocketbase-types-generator) (incomplete, unmaintained)
- [typed-pocketbase](https://github.com/david-plugge/typed-pocketbase) (incomplete, unmaintained)
- [pocketbase-schema](https://github.com/odama626/pocketbase-schema) (incomplete, unmaintained)

Comparison with [pocketbase-typegen](https://github.com/patmood/pocketbase-typegen)

| Aspect                          | **pb-typescript (this project)**                                 | **pocketbase-typegen**                  |
| ------------------------------- | ---------------------------------------------------------------- | --------------------------------------- |
| Integration model               | Native Go plugin, runs inside PocketBase                         | External generator                      |
| Requires PocketBase credentials | ❌ No                                                           | ✅ Yes                                   |
| Computed / virtual fields       | ✅ Available with [pb-computedfields](https://github.com/davenh99/pb-computedfields)| ❌ Not supported      |
| Select options export           | ✅ Optional generated constants                                   | ❌ Not provided                         |
| Completeness                    | ⚠️ Incomplete                                               | ✅ Comprehensive                    |
| Intended audience               | Developers writing custom pocketbase instances, extending with go           | Any pocketbase developer |

SO, why pb-typescript (this project) vs pocketbase-typegen?
If you are not extending with go, you can't use this package.

If you are extending with go, pb-typescript is simpler to setup (no requirement for setting up credentials, pb custom hooks).
pb-typescript includes support for computed fields, and including them in your typed output.
pb-typescript optionally exports select options as constants.

Some features missing in pb-typescript, but finished in pocketbase-typegen:
- full validation metadata
- richer JSON field schemas
- all field types not finished (missing geo point at least, haven't had a use for it yet)

These features will be added in future when needed.
