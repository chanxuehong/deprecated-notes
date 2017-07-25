# golang 数据库操作注意事项

## 不要忘记 defer rows.Close()

#### 正规的用法
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(2)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		list = append(list, id)
	}
	return list, rows.Err()
}

func main() {
	fmt.Println(test())
	fmt.Println(test())
	fmt.Println(test())
}
```

```Text
[107147279 107147287 107147294 107147295 107147304] <nil>
[107147279 107147287 107147294 107147295 107147304] <nil>
[107147279 107147287 107147294 107147295 107147304] <nil>
```

#### 虽然没有 defer rows.Close(), 但是也能工作, 因为 rows.Next() 最后一次会调用 rows.Close()
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(2)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		list = append(list, id)
	}
	return list, rows.Err()
}

func main() {
	fmt.Println(test())
	fmt.Println(test())
	fmt.Println(test())
}
```

```Text
[107147279 107147287 107147294 107147295 107147304] <nil>
[107147279 107147287 107147294 107147295 107147304] <nil>
[107147279 107147287 107147294 107147295 107147304] <nil>
```

#### 这种情况就不能工作了, 数据库链接没有得到释放, 会死锁
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(2)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		list = append(list, id)
		if len(list) >= 4 {
			break // 注意这里
		}
	}
	return list, rows.Err()
}

func main() {
	fmt.Println(test())
	fmt.Println(test())
	fmt.Println(test())
}
```

```Text
[107147279 107147287 107147294 107147295] <nil>
[107147279 107147287 107147294 107147295] <nil>
--------卡住
```

#### 上面代码如果用了 defer rows.Close() 就能工作
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(2)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() ([]int64, error) {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		list = append(list, id)
		if len(list) >= 4 {
			break
		}
	}
	return list, rows.Err()
}

func main() {
	fmt.Println(test())
	fmt.Println(test())
	fmt.Println(test())
}
```

```Text
[107147279 107147287 107147294 107147295] <nil>
[107147279 107147287 107147294 107147295] <nil>
[107147279 107147287 107147294 107147295] <nil>
```

## 数据库操作不要嵌套

#### 下面的例子虽然能够工作, 但是是有问题的
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(100)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	fmt.Println("goroutine started")
	time.Sleep(time.Second)

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			fmt.Println(err)
			return
		}
		list = append(list, id)
		if err = DB.Ping(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Wait()
	fmt.Println("done")
}
```

```Text
goroutine started
goroutine started
goroutine started
[107147279 107147287 107147294 107147295 107147304]
[107147279 107147287 107147294 107147295 107147304]
[107147279 107147287 107147294 107147295 107147304]
done
```

#### 上面的例子修改 DB.SetMaxOpenConns(3), 结果就不行了
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(3)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() {
	rows, err := DB.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	fmt.Println("goroutine started")
	time.Sleep(time.Second)

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			fmt.Println(err)
			return
		}
		list = append(list, id)
		if err = DB.Ping(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Wait()
	fmt.Println("done")
}
```

```Text
goroutine started
goroutine started
goroutine started
--------卡住
```

#### 同样对于事务也是一样的
```golang
package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const DSN = "root:123456@tcp(127.0.0.1:3306)/swift_nuochou_com?clientFoundRows=false&parseTime=true&loc=Asia%2FShanghai&timeout=5s&collation=utf8mb4_bin&interpolateParams=true"

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(3)
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}

func test() {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := tx.Query("SELECT id FROM users LIMIT 5")
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	defer rows.Close()

	fmt.Println("goroutine started")
	time.Sleep(time.Second)

	list := make([]int64, 0, 16)

	var id int64
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			fmt.Println(err)
			tx.Rollback()
			return
		}
		list = append(list, id)
		if err = DB.Ping(); err != nil {
			fmt.Println(err)
			tx.Rollback()
			return
		}
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	tx.Commit()
	fmt.Println(list)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		test()
	}()
	wg.Wait()
	fmt.Println("done")
}
```

```Text
goroutine started
goroutine started
goroutine started
--------卡住
```