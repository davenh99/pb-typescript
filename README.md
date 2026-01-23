## TypeScript types generator for pocketbase

### Usage

In your pocketbase project, run:

```
go get github.com/davenh99/pb-typescript
```

Then, in your main.go, put in something like (you can use an env var to avoid attaching the plugin in production):

```go
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

Here is a sample output (the preview is a computed field, attached using [pb-computedfields](https://github.com/davenh99/pb-computedfields):

```ts
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

You also have access to a typed pocketbase instance in pocketbase-types:

```ts
import { Collections, TypedPocketBase } from "../pocketbase-types";

const pb = new PocketBase("http://127.0.0.1:8090") as TypedPocketBase;

// exercise response here is typed as ExercisesRecord
const exercise = await pb.collection(Collections.Exercises).getFirstListItem(pb.filter("name = {:name}", { name }));
```

Because of the complexity required to properly type expands, if you are using expands, it would be recommended to write a new types file manually, and extend types there (possibly extended.d.ts or whatever you want to call it).
```ts
interface SessionsRecordExpand extends SessionsRecord {
  expand: {
    tags: TagsRecord[];
    sessionExercises_via_session: SessionExercisesRecordExpand[];
    sessionMeals_via_session: SessionMealsRecordExpand[];
  };
}

// you then have to type your requests manually as well.
import { Collections } from "../pocketbase-types";

const session = await pb.collection<SessionsRecordExpand>(Collections.Sessions).GetOne(id, {expand: "tags, sessionExercises_via_session, sessionMeals_via_session"})
```

For more examples of how the types look when they are generated check [progressa](https://github.com/davenh99/progressa/blob/main/ui/base.d.ts)

Comparison with [pocketbase-typegen](https://github.com/patmood/pocketbase-typegen)

| Aspect                          | **pb-typescript (this project)**                                 | **pocketbase-typegen**                  |
| ------------------------------- | ---------------------------------------------------------------- | --------------------------------------- |
| Integration model               | Native Go plugin, runs inside PocketBase                         | External generator                      |
| Requires PocketBase credentials | ❌ No                                                           | ✅ Yes                                   |
| Computed / virtual fields       | ✅ Available with [pb-computedfields](https://github.com/davenh99/pb-computedfields)| ❌ Not supported      |
| Completeness                    | ⚠️ Incomplete                                               | ✅ Comprehensive                    |
| Intended audience               | Developers writing custom pocketbase instances, extending with go           | Any pocketbase developer |

SO, why pb-typescript (this project) vs pocketbase-typegen?
If you are not extending with go, you can't use this package.

If you are extending with go, pb-typescript is simpler to setup (no requirement for setting up credentials, pb custom hooks).
pb-typescript includes support for computed fields, and including them in your typed output.
pb-typescript optionally exports select options as constants.

Some features missing in pb-typescript, but finished in pocketbase-typegen:
- full validation metadata
- all field types not finished (missing geo point at least, haven't had a use for it yet)

These features will be added in future when needed.
