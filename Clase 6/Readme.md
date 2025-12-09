# 5. Rust

## 1.5.1. Instalación de Rust

Siempre es bueno seguir lo que la [documentacion oficial de Rust dicta](https://www.rust-lang.org/tools/install). Pero en este caso la instalación de Rust para Linux es muy sencilla y únicamente necesitamos estos comandos.

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

Este comando aparte de instalar Rust, instala Rustup que es como el administrador no solo del lenguaje sino de otras herramientas, como las del compilado, documentación y manejo de crates.

Para verificar que se haya instalado todo correctamente vamos a correr los siguientes comandos.

```bash
rustc --version
```

Si por alguna razón no llegara a funcionar este comando, seguramente es que no se agregaron al PATH. Para gestionar ese tema les recomendamos seguir la [documentación de instalación del PATH](https://doc.rust-lang.org/book/ch01-01-installation.html)

## 1.5.2. Conceptos básicos de Rust

### Introducción a Rust

- **Historia y evolución**
  - Rust fue desarrollado por Mozilla Research, con el primer lanzamiento estable en 2015.
  - Se diseñó para abordar problemas de memoria y concurrencia que son comunes en lenguajes como C y C++.
- **Propósito y filosofía del lenguaje**
  - Rust se centra en la seguridad, velocidad y concurrencia.
  - Busca ofrecer un sistema de tipos que prevenga errores de memoria y condiciones de carrera sin necesidad de un recolector de basura.

### Inicio de un proyecto en Rust

Existen dos maneras de poder compilar código en Rust, la más sencilla es crear un archivo con la extensión `.rs` y luego ejecutar el comando `rustc main.rs`.

La segunda manera y la recomendada es utilizar **Cargo**, cargo es el manejador de paquetes en Rust, muchos de los **Rustaceans** usan esta herramienta para manejar sus paquetes propios como de terceros. Usualmente cargo ya viene instalado, para validar debemos escribir esto en la consola:

```bash
cargo --version
```

**Para iniciar un proyecto de rust solo basta con escribir el siguiente comando**

```bash
cargo new hello_cargo
cd hello_cargo
```

La estructura que maneja cargo es la siguiente:

```bash
hello_cargo
├── Cargo.toml
└── src
    └── main.rs
```

`Cargo.toml` es el archivo de configuración de nuestro proyecto, en el se especifican las dependencias y la información del proyecto.

`src/main.rs` es el archivo principal de nuestro proyecto, en el se escribe el código de nuestro programa.

**Con esto en mente**, ya podemos comenzar con conceptos básicos de Rust, pero antes, probemos hacer un hola mundo. Al hacer un proyecto nuevo con el comando anterior por lo general el código de hola mundo se crea por defecto, en caso de que no esté la función de hola mundo copiamos el siguiente código al .rs:

```rust
// fn define una función en Rust
fn main() {
    // println! es una macro (indicado por el !) que imprime a la consola
    // Las macros se ejecutan en tiempo de compilación
    println!("Hello, world!");
}
```

Finalmente ejecutamos el siguiente comando para compilar y ejecutar la aplicación:

```bash
cargo run
```

### Sintaxis básica

- **Variables y constantes**
  - Las variables en Rust son inmutables por defecto y se declaran con `let`.
  - Podemos hacer una variable mutable con `mut`.
  - Las constantes se declaran con `const` y deben tener un tipo explícito.

    ```rust
    // Variable inmutable - no podemos cambiar su valor después de asignarlo
    let x = 5;
    
    // Variable mutable - podemos modificar su valor
    let mut y = 10;
    y = 15; // Esto es válido porque y es mutable
    
    // Constante - debe tener tipo explícito y nunca puede cambiar
    // Por convención, las constantes se escriben en MAYÚSCULAS
    const MAX_POINTS: u32 = 100_000;
    
    // Ejemplo de shadowing - podemos redeclarar una variable
    let x = x + 1; // Ahora x vale 6
    let x = x * 2; // Ahora x vale 12
    ```

- **Tipos de datos primitivos**: Rust maneja una gran variedad de primitivos de los cuales son:
- **Enteros con signo (signed)**: i8, i16, i32, i64, i128, isize
- **Enteros sin signo (unsigned)**: u8, u16, u32, u64, u128, usize
- **Punto flotante**: f32, f64
- **Booleano**: bool
- **Carácter**: char
- **Compuestos**: tuplas y arreglos

    ```rust
    // Enteros con signo - pueden ser negativos
    let numero_negativo: i32 = -42;
    
    // Enteros sin signo - solo valores positivos
    let edad: u8 = 25;
    
    // Punto flotante - números decimales
    let pi: f64 = 3.14159;
    
    // Booleano - verdadero o falso
    let es_activo: bool = true;
    
    // Carácter - un solo carácter Unicode
    let letra: char = 'R';
    let emoji: char = '😊';
    
    // Tupla - puede contener diferentes tipos de datos
    let tup: (i32, f64, u8) = (500, 6.4, 1);
    // Podemos acceder a los elementos por índice
    let quinientos = tup.0;
    let seis_punto_cuatro = tup.1;
    
    // Arreglo - todos los elementos deben ser del mismo tipo
    let arr = [1, 2, 3, 4, 5];
    let primer_elemento = arr[0]; // Accedemos por índice
    ```

- **Arrays y Slices**: Los arrays son de tamaño fijo y los slices son una vista de un array.

    ```rust
    // Array de tamaño fijo - 5 elementos de tipo i32
    let arr: [i32; 5] = [1, 2, 3, 4, 5];
    
    // Crear un array con el mismo valor repetido
    let arr_repetido = [3; 5]; // [3, 3, 3, 3, 3]
    
    // Slice - referencia a una porción del array
    // El rango 1..3 incluye el índice 1 pero excluye el 3
    let slice = &arr[1..3]; // Contiene [2, 3]
    
    // Slice de todo el array
    let slice_completo = &arr[..];
    ```

- **Strings**: Rust maneja strings de dos formas: String y &str. Los strings son UTF-8 válidos y se pueden crear de varias formas. Mientras que &str es un slice de un string. Los String son mutables y los &str son inmutables. Eso quiere decir que los String se guardan en el heap y los &str en el stack.

    ```rust
    // String - tipo de dato que se almacena en el heap, mutable
    let s1 = String::from("hello");
    
    // &str - string literal, inmutable, se almacena en el stack
    let s2 = "world";
    
    // Concatenar strings usando la macro format!
    let s3 = format!("{} {}", s1, s2); // "hello world"
    
    // Crear un slice de string (&str) desde un String
    let s4: &str = &s1;
    
    // Añadir texto a un String mutable
    let mut s5 = String::from("Hola");
    s5.push_str(" Rust"); // "Hola Rust"
    s5.push('!'); // "Hola Rust!"
    ```

- **Structs**: Los structs son tipos de datos personalizados que permiten agrupar datos de diferentes tipos.

    ```rust
    // Definición de un struct - agrupa datos relacionados
    struct User {
        username: String,
        email: String,
        sign_in_count: u64,
        active: bool,
    }

    fn main() {
        // Crear una instancia del struct
        let user1 = User {
            username: String::from("user1"),
            email: String::from("prueba@gmail.com"),
            sign_in_count: 1,
            active: true,
        };
        
        // Acceder a los campos del struct
        println!("Usuario: {}", user1.username);
        
        // Struct mutable - podemos modificar sus campos
        let mut user2 = User {
            username: String::from("user2"),
            email: String::from("test@gmail.com"),
            sign_in_count: 0,
            active: false,
        };
        user2.sign_in_count += 1; // Incrementamos el contador
    }
    ```

- **Enums**: Los enums permiten definir un tipo que puede ser uno de varios valores. Otra forma de entender los enums son las variantes que puede llegar a tener un tipo.

    ```rust
    // Enum simple - define un tipo con variantes posibles
    enum IpAddrKind {
        V4,
        V6,
    }

    // Struct que usa el enum
    struct IpAddr {
        kind: IpAddrKind,
        address: String,
    }

    fn main() {
        // Crear instancias de las variantes del enum
        let four = IpAddrKind::V4;
        let six = IpAddrKind::V6;

        // Crear structs usando las variantes
        let home = IpAddr {
            kind: IpAddrKind::V4,
            address: String::from("127.0.0.1"),
        };

        let loopback = IpAddr {
            kind: IpAddrKind::V6,
            address: String::from("::1"),
        };
        
        // Enum con datos asociados - más idiomático en Rust
        enum IpAddrMejor {
            V4(String),
            V6(String),
        }
        
        let home_mejor = IpAddrMejor::V4(String::from("127.0.0.1"));
    }
    ```

### Control de flujo

- **Condicionales**
   - Uso de if, else if, y else

    ```rust
    let number = 6;
    
    // Estructura if-else if-else estándar
    if number % 4 == 0 {
        println!("number is divisible by 4");
    } else if number % 3 == 0 {
        println!("number is divisible by 3");
    } else {
        println!("number is not divisible by 4, 3, or 2");
    }
    
    // if es una expresión, puede devolver un valor
    let resultado = if number > 5 {
        "mayor que 5"
    } else {
        "menor o igual a 5"
    };
    
    // Los tipos deben coincidir en todas las ramas
    let numero = if number > 5 { 10 } else { 20 };
    ```

- **Bucles**
   - loop: Ejecuta un bloque de código repetidamente hasta que se interrumpa con break.
   - while: Ejecuta un bloque de código mientras una condición sea verdadera.
   - for: Itera sobre una colección de elementos.

    ```rust
    // loop - ciclo infinito que se detiene con break
    let mut count = 0;
    loop {
        count += 1;
        println!("Count: {}", count);
        
        // break termina el bucle
        if count == 3 {
            break;
        }
    }
    
    // loop puede devolver un valor con break
    let resultado = loop {
        count += 1;
        if count == 10 {
            break count * 2; // Devuelve 20
        }
    };

    // while - ejecuta mientras la condición sea verdadera
    while count != 0 {
        println!("Countdown: {}", count);
        count -= 1;
    }
    
    // for - itera sobre un rango o colección
    // 1..4 es un rango que va de 1 a 3 (no incluye el 4)
    for number in 1..4 {
        println!("{}", number);
    }
    
    // Iterar sobre un array
    let arr = [10, 20, 30, 40, 50];
    for elemento in arr.iter() {
        println!("El valor es: {}", elemento);
    }
    
    // Iterar con índice
    for (indice, valor) in arr.iter().enumerate() {
        println!("Índice {}: valor {}", indice, valor);
    }
    ```

### Funciones

- **Definición y llamada a funciones**
  - Las funciones se definen con fn.
  - En Rust manejamos el concepto de snake_case para nombrar funciones y variables.

    ```rust
    fn main() {
        // Llamamos a la función
        say_hello();
        
        // Llamamos con parámetros
        say_hello_to("Rust");
    }

    // Función simple sin parámetros ni retorno
    fn say_hello() {
        println!("Hello!");
    }
    
    // Función con parámetro
    // El tipo del parámetro debe especificarse
    fn say_hello_to(name: &str) {
        println!("Hello, {}!", name);
    }
    ```

- **Parámetros y retorno de valores**
   - Las funciones pueden recibir parámetros y devolver valores.
   - El tipo de retorno se especifica con ->.
   - Muchas funciones en Rust pueden o no manejar la palabra reservada return. Con Rust podemos omitir el return.

    ```rust
    fn main() {
        let result = add(5, 3);
        println!("Result: {}", result); // Imprime: Result: 8
        
        // Llamada con diferentes valores
        let suma = add(10, 20);
        println!("Suma: {}", suma); // Imprime: Suma: 30
    }

    // Función con parámetros y valor de retorno
    // -> i32 indica que devuelve un entero de 32 bits
    fn add(x: i32, y: i32) -> i32 {
        // La última expresión sin punto y coma es el valor de retorno
        x + y
    }
    
    // También podemos usar return explícitamente
    fn subtract(x: i32, y: i32) -> i32 {
        return x - y; // return es opcional pero válido
    }
    
    // Función con múltiples parámetros y lógica condicional
    fn max(a: i32, b: i32) -> i32 {
        // if como expresión que retorna un valor
        if a > b {
            a
        } else {
            b
        }
    }
    ```

## 1.5.3 Conceptos avanzados

### Ownership

Rust es un lenguaje muy especial ya que cuenta con características únicas como es el Ownership. Pero antes de explicar qué es el Ownership, es importante entender cómo otros lenguajes manejan la memoria.

- **Manejo de memoria en otros lenguajes**
  - **C**: En C se maneja la memoria de forma manual, debemos liberar la memoria manualmente.
  - **C++**: En C++ se maneja la memoria de forma manual, pero podemos usar punteros inteligentes.
  - **Java, C#, Python**: Estos lenguajes cuentan con un recolector de basura que se encarga de liberar la memoria.

Una forma sencilla de entender cómo el Garbage Collector funciona es que se encarga de liberar la memoria que ya no se necesita. Pero en Rust no contamos con un recolector de basura, en su lugar contamos con Ownership.

**Entiendiendo Ownership**

Con el uso de este concepto, hace que Rust sea un lenguaje seguro y rápido sin la necesidad de un GC.

**Reglas de Ownership**
- Cada valor en Rust tiene una variable que es su dueño.
- Solo puede haber un dueño a la vez.
- Cuando el dueño sale del alcance, el valor se libera.

Con esto en mente vamos a ver un ejemplo de cómo se maneja la memoria en Rust.

```rust
fn main() {
    {                    // s no es válido aquí, no está declarado todavía
        let s = String::from("hello");   // s es válido desde este punto
        
        // Podemos usar s aquí
        println!("{}", s);
        
    }                   // s ya no es válido, sale del alcance
                       // Rust llama automáticamente a drop() para liberar memoria

    // println!("{}", s); // ERROR: s no existe en este alcance
}
```

Cosas a tomar en cuenta acá es que cuando s sale del alcance, Rust se encarga de liberar la memoria. Esto se hace con un concepto llamado Drop, que se encarga de liberar la memoria. Rust lo que hace en esta parte es que una vez deja de ser válida, se libera memoria.

Por ejemplo, si usualmente cuando vemos esto

```rust
// Tipos primitivos en el stack - se copian
let x = 5;
let y = x; // Se copia el valor, ambos son válidos

println!("x = {}, y = {}", x, y); // Ambos funcionan
```

Lo que muchos lenguajes y Rust hacen es que se copia el valor de x a y. Esto es posible porque estamos hablando de un tipo primitivo estático que se encuentra en el stack. Pero si hablamos de un tipo compuesto como un String, Rust no copia el valor, sino que mueve la referencia a la memoria.

```rust
// Tipos en el heap - se mueven (move)
let s1 = String::from("hello");
let s2 = s1; // s1 se mueve a s2, s1 ya no es válido

// println!("{}", s1); // ERROR: s1 ya no es válido
println!("{}", s2); // Solo s2 es válido ahora

// Si queremos copiar, usamos clone()
let s3 = String::from("world");
let s4 = s3.clone(); // Se crea una copia profunda

println!("s3 = {}, s4 = {}", s3, s4); // Ambos son válidos
```

En esta figura, podemos ver lo que realmente Rust hace, en el Stack almacenamos información de S1 como el puntero, la longitud y la capacidad. Y en el Heap almacenamos el valor de "hello". 
![Ownership](./img/ownership.png)

Entonces cuando nosotros hacemos el `let s2 = s1`, s2 ahora apunta al mismo puntero que S1.

![Ownership](./img/ownership2.png)

Pero si nos damos cuenta esto va a provocar un error, ya que tendríamos lo que se conoce como un double free. **Para evitar esto Rust invalida a S1, es decir, S1 ya no es válida.**

Rust lo que nos dice ahora es que se realizó un movimiento de la memoria, es decir, S1 ya no es válida y S2 es la que tiene la referencia a la memoria.

**Funciones y Ownership**

El mecanismo es muy similar a lo que vimos anteriormente, cuando pasamos un valor a una función, Rust mueve la referencia a la memoria.

```rust
fn main() {
    let s = String::from("hello");  // s entra en el ámbito

    takes_ownership(s);             // s se mueve dentro de la función
                                    // s ya no es válido aquí
    
    // println!("{}", s); // ERROR: s fue movido

    let x = 5;                      // x entra en el ámbito

    makes_copy(x);                  // x se copia dentro de la función
                                    // x sigue siendo válido aquí porque i32 es Copy
    
    println!("{}", x); // Funciona: x todavía es válido

} // Aquí, x sale del alcance, luego s. Pero como el valor de s fue movido, no pasa nada especial.

fn takes_ownership(some_string: String) { // some_string entra en el ámbito
    println!("{}", some_string);
} // Aquí, some_string sale del alcance y `drop` es llamado. El espacio en memoria es liberado.

fn makes_copy(some_integer: i32) { // some_integer entra en el ámbito
    println!("{}", some_integer);
} // Aquí, some_integer sale del alcance. Nada especial sucede porque es Copy.
```

**Retorno de valores**

Cuando una función retorna un valor, Rust también mueve la referencia a la memoria.

```rust
fn main() {
    // gives_ownership mueve su valor de retorno a s1
    let s1 = gives_ownership();

    let s2 = String::from("hello");     // s2 entra en el ámbito

    // s2 se mueve a la función, que también retorna su valor a s3
    let s3 = takes_and_gives_back(s2);
    
    println!("s1 = {}", s1); // Válido
    // println!("{}", s2); // ERROR: s2 fue movido
    println!("s3 = {}", s3); // Válido
}

// Esta función crea un String y lo retorna
fn gives_ownership() -> String {
    let some_string = String::from("yours"); // some_string entra en el ámbito

    some_string  // some_string se retorna y se mueve a la función que llama
}

// Esta función toma un String y retorna uno
fn takes_and_gives_back(a_string: String) -> String { // a_string entra en el ámbito

    a_string  // a_string se retorna y se mueve a la función que llama
}
```

**Referencias y Borrowing**

Para poder evitar que Rust mueva la referencia a la memoria, podemos utilizar lo que se conoce como referencias. Las referencias permiten que múltiples partes del código accedan a los mismos datos sin necesidad de mover la referencia.

En el ejemplo anterior vimos que para obtener el tamaño de un String se tenía que regresar su tamaño y luego el String. Pero con las referencias podemos hacer que la función no mueva la referencia a la memoria.

```rust
fn main() {
    let s1 = String::from("hello");

    // &s1 crea una referencia a s1 sin mover ownership
    let len = calculate_length(&s1);

    // s1 todavía es válido porque solo prestamos una referencia
    println!("The length of '{}' is {}.", s1, len);
}

// s es una referencia a un String (no toma ownership)
fn calculate_length(s: &String) -> usize {
    s.len()
} // Aquí, s sale del alcance pero como no tiene ownership, no se libera nada
```

**Referencias mutables**

Las referencias mutables permiten modificar el valor de un dato.

```rust
fn main() {
    // La variable debe ser mutable
    let mut s = String::from("hello");

    // &mut s crea una referencia mutable
    change(&mut s);
    
    println!("{}", s); // Imprime "hello, world"
}

// some_string es una referencia mutable, puede modificar el String
fn change(some_string: &mut String) {
    // Agregamos texto al String
    some_string.push_str(", world");
}

// Ejemplo de restricción: solo una referencia mutable a la vez
fn ejemplo_restriccion() {
    let mut s = String::from("hello");
    
    let r1 = &mut s;
    // let r2 = &mut s; // ERROR: no podemos tener dos referencias mutables
    
    println!("{}", r1);
}
```

**Reglas de Referencias**   
- Solo podemos tener una referencia mutable a la vez.
- No podemos tener una referencia mutable y una inmutable al mismo tiempo.
- Las referencias deben estar dentro del alcance.

```rust
fn main() {
    let mut s = String::from("hello");
    
    // Múltiples referencias inmutables están permitidas
    let r1 = &s;
    let r2 = &s;
    println!("{} and {}", r1, r2);
    // r1 y r2 ya no se usan después de este punto
    
    // Ahora podemos crear una referencia mutable
    let r3 = &mut s;
    r3.push_str(" world");
    println!("{}", r3);
}
```

### Implementaciones en Rust

Las implementaciones en Rust permiten agregar métodos a un struct o enum.

```rust
// Definimos un struct
struct Rectangle {
    width: u32,
    height: u32,
}

// impl define un bloque de implementación para Rectangle
impl Rectangle {
    // Método que toma &self (referencia inmutable a la instancia)
    fn area(&self) -> u32 {
        // Calculamos el área
        self.width * self.height
    }

    // Método que compara dos rectángulos
    fn can_hold(&self, other: &Rectangle) -> bool {
        // Verifica si este rectángulo puede contener a otro
        self.width > other.width && self.height > other.height
    }
    
    // Función asociada (no toma self) - constructor personalizado
    fn square(size: u32) -> Rectangle {
        Rectangle {
            width: size,
            height: size,
        }
    }
}

fn main() {
    let rect1 = Rectangle { width: 30, height: 50 };
    let rect2 = Rectangle { width: 10, height: 40 };
    let rect3 = Rectangle { width: 60, height: 45 };

    // Llamamos al método area usando la sintaxis de punto
    println!("The area of the rectangle is {} square pixels.", rect1.area());

    // Llamamos al método can_hold
    println!("Can rect1 hold rect2? {}", rect1.can_hold(&rect2));
    println!("Can rect1 hold rect3? {}", rect1.can_hold(&rect3));
    
    // Llamamos a la función asociada usando ::
    let sq = Rectangle::square(20);
    println!("Square area: {}", sq.area());
}
```

### Traits y generics

**Traits**

Los traits son una forma de definir comportamientos en Rust. Los traits permiten definir métodos que un tipo debe implementar.

```rust
// Definimos un trait - un conjunto de métodos que los tipos pueden implementar
pub trait Summary {
    fn summarize(&self) -> String;
    
    // Método con implementación por defecto
    fn default_summary(&self) -> String {
        String::from("(Read more...)")
    }
}

pub struct NewsArticle {
    pub headline: String,
    pub location: String,
    pub author: String,
    pub content: String,
}

// Implementamos el trait Summary para NewsArticle
impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {} ({})", self.headline, self.author, self.location)
    }
}

pub struct Tweet {
    pub username: String,
    pub content: String,
}

// Implementamos el mismo trait para otro tipo
impl Summary for Tweet {
    fn summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }
}

fn main() {
    let article = NewsArticle {
        headline: String::from("Rust 1.50 Released"),
        location: String::from("San Francisco"),
        author: String::from("Jane Doe"),
        content: String::from("The Rust team is happy to announce..."),
    };
    
    // Usamos el método del trait
    println!("New article: {}", article.summarize());
    
    let tweet = Tweet {
        username: String::from("rustlang"),
        content: String::from("Rust is awesome!"),
    };
    
    println!("New tweet: {}", tweet.summarize());
}
```

**Generics**

Los generics permiten definir funciones, structs y enums que pueden trabajar con diferentes tipos de datos.

```rust
// Struct genérico - T es un tipo que se define al crear la instancia
pub struct Point<T> {
    x: T,
    y: T,
}

// Implementación para cualquier tipo T
impl<T> Point<T> {
    // Método que retorna una referencia al campo x
    pub fn x(&self) -> &T {
        &self.x
    }
}

// Implementación específica para Point<f64>
impl Point<f64> {
    // Este método solo está disponible para Point<f64>
    pub fn distance_from_origin(&self) -> f64 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

// Función genérica que funciona con cualquier tipo T que implemente PartialOrd
fn largest<T: PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];
    
    for item in list {
        if item > largest {
            largest = item;
        }
    }
    
    largest
}

fn main() {
    // Point con enteros
    let p = Point { x: 5, y: 10 };
    println!("p.x = {}", p.x());

    // Point con flotantes
    let p = Point { x: 1.0, y: 4.0 };
    println!("p.x = {}", p.x());
    println!("Distance: {}", p.distance_from_origin());
    
    // Función genérica con diferentes tipos
    let numbers = vec![34, 50, 25, 100, 65];
    println!("The largest number is {}", largest(&numbers));
    
    let chars = vec!['y', 'm', 'a', 'q'];
    println!("The largest char is {}", largest(&chars));
}
```

**Derivación automática de Traits**

Rust permite derivar automáticamente algunos traits como Debug y Clone.

```rust
// #[derive] genera automáticamente la implementación de traits
#[derive(Debug, Clone, PartialEq)]
struct Rectangle {
    width: u32,
    height: u32,
}

fn main() {
    let rect = Rectangle { width: 30, height: 50 };
    
    // Debug nos permite usar {:?} para imprimir
    println!("rect is {:?}", rect);
    
    // Clone nos permite crear una copia
    let