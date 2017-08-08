# Go文本类型:string,byte,rune

## 概览

Go中，字符串`string`是内置类型，与文本处理相关的内置类型还有符文`rune`和字节`byte`。

与C语言类似，大多数关于字符串处理的函数都放在标准库里。Go将大部分字符串处理的函数放在了`strings`,`bytes`这两个包里。因为在字符串和整型间没有隐式类型转换，字符串和其他基本类型的转换的功能主要在标准库`strconv`中提供。`unicode`相关功能在`unicode`包中提供。

* Go语言源代码总是采用`UTF-8`编码
* 字符串`string`可以包含任意字节序列，通常是`UTF-8`编码的。
* 字符串字面值，在不带有字节转义的情况下一定是`UTF-8`编码的。
* Go使用`rune`代表`Unicode`**码位**。一个**字符**可能由一个或多个码位组成（复合字符）
* Go string是建立在**字节数组**的基础上的，因此对string使用`[]`索引会得到字节`byte`而不是字符`rune`。
* Go语言的字符串不是正规化(normalized)的，因此同一个字符可能由不同的字节序列表示。使用`unicode/norm`解决此类问题。



## 基础数据结构

### 数组与切片

数组`Array`是固定长度的数据结构，不存放任何额外的信息。通常不直接使用，而是作为切片的底层存储。

切片`Slice`描述了数组中一个连续的片段。在底层实现中，可以看成一个由三个word组成的结构体，这里word是CPU的字长。这三个字分别是`ptr`,`len`,`cap`，分别代表数组首元素地址，切片的长度，当前切片头位置到底层数组尾部的距离。

![](./media/godata3.png)





### 字节`byte`

字节`byte`实际上就是`uint8`的别名，只是为了和其他8bit类型相区别才单独起了别名。通常出现的更多的是字节切片`[]byte`与字节数组`[...]byte`。

字节可以用单引号扩起的单个字符表示，不过这种字面值和rune的字面值很容易搞混。赋予字节变量一个超出范围的值，如果在编译期能检查出来就会报`overflows byte`编译错误。

对于字节数组`[]byte`，实质上可以看做`[]uint8`，所以字节数组的本体结构定义如下：

```go
type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```





### 字符串`string`

字符串通常是UTF8编码的文本，由一系列8bit字节组成。`raw string literal`和不含转义符号的`string literal`一定是`UTF-8`编码的，但`string`其实可以含有任意的字节序列。

字符串是不可变对象，可以空(`s=""`)，但不会是`nil`。

![](./media/godata2.png)

底层实现上字符串与切片头类似，都带有一个指针，一个长度，但因为字符串一旦创建就不可变所以不需要`cap`字段，所以字符串其实由两个Word组成。64位机器上就是16个字节。

```go
type StringHeader struct {
    Data uintptr
    Len  int
}
```

关于`string`，有这么几点需要注意：

1. `string`常量会在编译期分配到**只读段**，对应数据地址不可写入，相同的`string`常量不会重复存储。 
2. `fmt.Sprintf`生成的字符串分配在堆上，对应数据地址可修改。 
3. 常量空字符串有数据地址，动态生成的字符串没有设置数据地址 
4. Golang string和[]byte转换,会将数据复制到堆上，返回数据指向复制的数据 
5. 动态生成的字符串，即使内容一样，数据也是在不同的空间 
6. 只有动态生成的string，数据可以被黑科技修改 
7. string和[]byte通过复制转换，性能损失接近4倍




### 符文`rune`

符文`rune`其实是`int32`的别名，表示一个Unicode的**码位**。

注意一个**字符(Character)**可以由一个或多个**码位(Code Point)**构成。例如带音调的`e`，即`é`,既可以由`\u00e9`单个码位表示，也可以由`e`和口音符号`\u0301`复合而成。这涉及到[normalization](https://blog.golang.org/normalization)的问题。但通常情况下一个字符就是一个码位。

```python
>>> print u'\u00e9', u'e\u0301',u'e\u0301\u0301\u0301'
é é é́́
```

符文的字面值是用单引号括起的一个或多个字符，例如`a`,`啊`,`\a`,`\141`,`\x61`,`\u0061`,`\U00000061`，都是合法的rune literal。其格式定义如下：

```
rune_lit         = "'" ( unicode_value | byte_value ) "'" .
unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\` octal_digit octal_digit octal_digit .
hex_byte_value   = `\` "x" hex_digit hex_digit .
little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .
big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
                           hex_digit hex_digit hex_digit hex_digit .
escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .
```

其中，八进制的数字范围是0~255，Unicode转义字符通常要排除0x10FFFF以上的字符和surrogate字符。



## 文本类型转换

### `string`与`[]byte`的转换

#### 常规做法

通常`[]byte`和`string`可以直接通过类型名强制转化，但实质上执行了一次堆复制。理论上`stringHeader`只是比`sliceHeader`少一个`cap`字段，但因为`string`需要满足不可变的约束，而`[]byte`是可变的，因此在执行`[]byte`到`string`的操作时会进行一次复制，在堆上新分配一次内存。

```go
// byte to string
s := string(b)

// string index -> byte
s[i] = b

// []byte to string
s := string(bytes)

// string to []byte
bytes := []byte(s)
```

#### 黑魔法

利用`unsafe.Pointer`和`reflect`包可以实现很多禁忌的黑魔法，但这些操作对GC并不友好。最好不要尝试。

```go
type Bytes []byte

// 将string转换为[]byte，'可以修改'，很危险，因为[]byte结构要多一个cap字段。
func StringBytes(s string) Bytes {
    return *(*Bytes)(unsafe.Pointer(&s))
}

// 不拷贝地将[]byte转换为string
func BytesString(b []byte) String {
    // 因为[]byte的Header只比string的Header多一个Cap字段。可以直接强制成`*String` 
    return *(*String)(unsafe.Pointer(&b))
}

// 获取&s[0]，即存储字符串的字节数组的地址指针，Go里不允许这种操作。 
func StringPointer(s string) unsafe.Pointer {
    p := (*reflect.StringHeader)(unsafe.Pointer(&s))
    return unsafe.Pointer(p.Data)
}

// r获取&b[0]，即[]byte底层数组的地址指针，Go里不允许这种操作
func BytesPointer(b []byte) unsafe.Pointer {
    p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    return unsafe.Pointer(p.Data)
}
```



### `string`与`rune`的转换

`string`和`rune`之间也可以通过类型名直接转换，不过`string`不能直接转换成单个的`rune`。

```go
// rune to string
str := string(r)

// range string -> rune
for i,r := range str

// string to []rune
runes := []rune(str)

// []rune to string
str := string(runes)
```

Go对于`UTF-8`有特殊的支持和处理（因为`UTF-8`和`Go`都是`Ken`发明的……。），这体现在对于`string`的`range`迭代上。

```go
const nihongo = "日本語"
for index, runeValue := range nihongo {
    fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
}

U+65E5 '日' starts at byte position 0
U+672C '本' starts at byte position 3
U+8A9E '語' starts at byte position 6
```

直接索引`string`会得到字节序号和相应字节。而对`string`进行`range`迭代，获得的就是字符`rune`的索引与相应的`rune`。



### `byte`与`rune`的转换

`byte`其实是`uint8`，而`rune`实际就是`int32`，所以两者之间的转换就是整数的转换。两个不同类型的整形数组之间是没有强制转换的方法的，只有通过`string`来曲线救国：`runes := []rune(string(bytes))`