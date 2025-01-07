# SnowFlake

```go
    generator, err := NewSnowflake(1)
	if err != nil {
		panic(err)
	}
	
	id, err := generator.NewID()
	if err != nil {
		panic(err)
	}
	fmt.Println(id.Int64())
```

 