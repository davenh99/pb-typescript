## TypeScript types generator for pocketbase

### Usage

In your pocketbase project, run:

```
go get github.com/davenh99/pb-typescript
```

Then, in your main.go, put in something like:

```
gentypes.Register(app, gentypes.Config{
    FilePath:                   "ui",
    PrintSelectOptions:         true,
})
```

Your types will then be generated in the ui directory, whenever the schema changes or you run (assuming your app executable is named pocketbase):

```
./pocketbase gen-types
```

for more examples of how the types look when they are generated check github.com.davenh99/progressa
