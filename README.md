# go-cache

## 简介

> hyper performance local cache written by Go！

**基于Go编写的，高性能本地缓存。**

特点：

- 高性能，千万级读性能，百万级写性能
- redis风格，可以像使用redis一样

## 使用方式
**获取包**
```go
go get github.com/wk331100/go-cache
```
## 调用示例
```go
func main() {
	// 初始化缓存
	c := go_cache.NewCache()
	// 写入String类型的缓存
	c.Set("name", "ZhangSan")
	c.Set("age", "18")
	// 读取字符类型的缓存
	if name, err := c.Get("name"); err == nil {
		fmt.Println(name)
	}

	u := &user{
		Uid:  1001,
		Name: "zhangSan",
		Age:  18,
	}
	bz, _ := json.Marshal(u)
	// 写入Hash类型的缓存
	c.HSet("users", "1001", bz)
	// 读取Hash类型的缓存
	if bzCache, err := c.HGet("users", "1001"); err == nil {
		u1 := &user{}
		fmt.Println(bzCache)
		if err = json.Unmarshal(bzCache.([]byte), u1); err == nil {
			fmt.Println(u1)
		}
	}
}

type user struct {
	Uid  int
	Name string
	Age  int
}
```


## 性能测试
### 性能汇总

| 数据结构  | 方法     | 操作 |    QPS |
|-------|--------|----|-------:|
| string | set()  | 写   |   93 万 |
| string | get()  | 读  | 1936 万 |
| hash   | hset() | 写  |  216 万 |
| hash   | hget() | 读  | 1728 万 |

### string 写性能: 93万QPS
```go
func TestSetStringQPS(t *testing.T) {
	start := time.Now()
	loop := 1000000
	for j := 0; j < loop; j++ {
		c.Set(strconv.Itoa(j), "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set string: %v\n", d.Seconds())
	fmt.Printf("benchmark set string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
```
结果：
```
=== RUN   TestSetStringQPS
benchmark set string: 1.064814931
benchmark set string qps: 939130 
--- PASS: TestSetStringQPS (1.06s)
PASS
```

### string 修改值性能: 2578万QPS
```go
func TestSetSameStringQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
    loop := 1000000
	for j := 0; j < loop; j++ {
		c.Set(key, "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set same string: %v\n", d.Seconds())
	fmt.Printf("benchmark set same string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
```
结果：
```
=== RUN   TestSetSameStringQPS
benchmark set same string: 0.038788867
benchmark set same string qps: 25780593 
--- PASS: TestSetSameStringQPS (0.04s)
PASS
```

### string 读性能: 1936万QPS
```go
func TestGetStringQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		_, _ = c.Get(key)
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark get string: %v\n", d.Seconds())
	fmt.Printf("benchmark get string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
```
结果：
```
=== RUN   TestGetStringQPS
benchmark get string: 0.051630968
benchmark get string qps: 19368221 
--- PASS: TestGetStringQPS (0.05s)
PASS
```

### hash 写性能: 216万QPS
```go
func TestSetHashQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		c.HSet(key, strconv.Itoa(j), "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set hash: %v s\n", d.Seconds())
	fmt.Printf("benchmark set hash qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
```
结果：
```
=== RUN   TestSetHashQPS
benchmark set hash: 0.462647734 s
benchmark set hash qps: 2161472 
--- PASS: TestSetHashQPS (0.46s)
PASS
```

### hash 读性能: 1728万QPS
```go
func TestGetHashQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		_, _ = c.HGet(key, strconv.Itoa(j))
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark get string: %v\n", d.Seconds())
	fmt.Printf("benchmark get string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
```
结果：
```
=== RUN   TestGetHashQPS
benchmark get string: 0.057859717
benchmark get string qps: 17283182 
--- PASS: TestGetHashQPS (0.06s)
PASS
```