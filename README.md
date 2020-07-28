# TIS 80
Tis80 es un pequeño y simple emulador de un *"fantasy computer"*. Un *"fantasy computer"* es, como su nombre indica, un ordenador ficticio, normalmente retro (de 8 bits, estilo Commodore64, por ejemplo). No es una [*"fantasy console"*](https://github.com/topics/fantasy-console) como el Pico8 o el Tic8. El objetivo de este emulador es solo emular el comportamiento una máquina de 8 bits. No es una consola de videojuegos.

## Roadmap
Aunque el proyecto ha cumplido el objetivo que me había propuesto, puede que se añadan las siguientes características:

* Cambio en la memoria de vídeo para permitir mostrar colores.
* Añadir sprites estilo Commodore64.
* Al pulsar la tecla F11 (o similar) el emulador se pausa y se muestra el estado actual de la CPU (memoria, registros y flags).
* Añadir soporte para operaciones aritmético-lógicas con enteros de 8 bits con signo.
* Añadir soporte para operaciones aritmético-lógicas con números de coma flotante (aún no se tiene clara su implementación).
* Añadir interrupciones para teclado (aún no se tiene clara su implementación).

## Estructura del proyecto
El proyecto se divide en:

* __Un assembler/diassembler__: Estos programas te van a servir para pasar tu código ensamblador a un binario que entiende Tis80. Ambos están hechos en Go y se encuentran en el directirio *"asm"*.

* __CPU__: Se trata del código principal del Tis80. Se encuentra bajo el directorio *"cpu"*.

* __Una implementación de consola__: Una versión para usar en la shell. Está en el directorio *"console"*.

* __Una implementación gráfica__: Por el momento muy limitada. Es exactamente igual que la versión de consola con la diferencia que en esta versión aparece una ventana de 320x200 píxeles donde se muestra el contienido de la memoria de vídeo del emulador. Se encuentra en *"screen"*.

## Compilación
Para compilar el proyecto necesitas:

* Go v1.14
* Un compilador de C. Lo he probado con GCC sobre GNU/Linux y sobre CLang en OSx
* Make
* Para la versión gráfica necesitarás también SDL2 y SDL2-TTF.

Teniendo esto, simplemente situate en la raíz del proyecto y ejecuta make. Si todo ha ido bien deberías de tener una carpeta llamda *build* con los siguientes binarios:

* __tisasm__: ensamblador
* __tisdiasm__: desensamblador
* __tisconsole__: versión del emulador para la línea de comandos
* __tis__: versión del emulador gráfica.

## Arquitectura del Tis80
El Tis80 se trata de un ordenador de 8 bits, con un rango de direcciones de 64K palabras. Dispone de 16 registros de uso general de 8 bits y de un acumulador (ACC) también de 8 bits. Las flags disponibles son:

* ACC_OVERFLOW [0]: Identificado como la flag 0. El acumulador ha tenido un overflow. Esto ocurre cuando:
    * El contenido del acumulador es 0xFFFF y se suma uno o más.
    * El contenido del acumulador es 0x0000 y se resta uno o más.
* STACK_OVERFLOW [1]: Identificado como la flag 1. Se ha intentado hacer un push en el stack cuando ya se ha alcanzado el límite.
* IO_ERROR [2]: Identificado como la flag 2. Se ha intentado leer/escribir en un floppy-disk y no ha sido posible.

Todos estos flags, cuando se activan, generan una interrupción.

Todas las operaciones aritmético-lógicas de la máquina son con enteros de 8 bits sin signo. Próximamente se permitirá realizar otras operaciones (mirar en Roadmap).

## Mapa de memoria
El Tis80 tiene un rango de direcciones de 64K. La palabra de memoria es de 8 bits. Teniendo esto en cuenta, el mapa de memoria del Tis80 es:

| Direcciones   |                                                                                                                                     Significado                                                                                                                                    |
|---------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
| $0000 - $00FF | Vectores de interrupciones. Se espera que en las siguientes direcciones se guarden el código de las siguientes interrupciones: * $0000: ACC_OVERFLOW * $0002: STACK_OVERFLOW * $0004: IO_ERROR * $0006: KEYBOARD El resto de vectores se dejan libre para funcionalidad del Kernel |
| $0100 - $0103 | Parámetros para subrutinas. Dependiendo del caso se puede usar de distintas formas                                                                                                                                                                                                 |
| $0104 - $01FF | Stack                                                                                                                                                                                                                                                                              |
| $0200 - $2FFF | Código del kernel (ROM)                                                                                                                                                                                                                                                            |
| $3000 - $3FFF | Memoria de vídeo                                                                                                                                                                                                                                                                   |
| $4000 - $40FF | Buffer de entrada por teclado                                                                                                                                                                                                                                                      |
| $4100 - $FFFF | RAM (memoria para los programas del usuario)                                                                                                                                                                                                                                       |

## Ensamblador del Tis80
### Tipos de datos

* __Secciones__: son un punto seguido de un texto. Solo hay dos: *.data* y *.code*.
* __Strings__: Son unas comillas dobles, seguidas de un texto y terminadas en unas comillas dobles. Solo pueden aparecer en la sección de datos. Por ejemplo: "Hola" o "Bienvenido al Tis80".
* __Números en hexadecimal__: Es cualquier número que empieza por un 0, seguido por una x, continuado por un número hexadecimal (es decir, se admite dígitos y las letras 'a', 'b', 'c', 'd', 'e' y 'f' tanto en minusculas como mayusculas).
* __Números en decimal__: Cualquier número del 1 al 255 (los números están limitados a 8 bits). Si se quiere escribir el 0 en decimal, se debe hacer usando la notación hexadecimal.
* __Tags__: Son equivalentes a las direcciones de memoria. Útiles para destinos de saltos. Se declaran con dos puntos (por ejemplo :destino). Se usan escribiendo el nombre de la tag sin los dos puntos (por ejemplo **jmp destino**). Se tranforman en direcciones fijas cuando se ensambla.
* __Comentarios__: Comienzan con punto y coma (;) y terminan al final de la línea.

### Secciones
Una sección es una parte del código ensamblador dedicada para indicar información de distinto tipo al emulador. Existen dos secciones:

* __Sección de datos__: Se declaran strings o números junto a la dirección de inicio de estos datos. Se escribe como *.data*
* __Sección de código__: Se declaran las instrucciones a ejecutar. Se escribe como *.code* y justo depués debe aparecer la dirección de memoria a partir de la cual va a ser escrito el código en la memoria. Es decir, si la sección de código se declara como *.code $0200* significa que a partir de la dirección $0200 (inclusive) se empezará a escribir el código.

### Conjunto de instrucciones

Aritmético-Lógicas
| Instrucción | OpCode |  Operación | Explicacón |
|-|-|:-:|-|
| add Rx | 0x01 | acc + Rx -> acc |  |
| addi INT | 0x02 | acc + INT -> acc |  |
| sub Rx | 0x03 | acc - Rx -> acc |  |
| subi INT | 0x04 | acc - INT -> acc |  |
| sil | 0x05 | acc << 1 -> acc |  |
| sir | 0x06 | acc >> 1 -> acc |  |
| and Rx | 0x07 | acc ^ Rx -> acc |  |
| or Rx | 0x08 | acc (or) Rx -> acc |  |
| not | 0x09 | ¬acc -> acc |  |
| xor Rx | 0x0a | acc (xor) Rx -> acc |  |

Salto
| Instrucción | OpCode | Operación | Explicacón |
|-|-|:-:|-|
| jmp MEM | 0x20 | pc = MEM | Salto incondicional a la dirección MEM. Por ejemplo jmp $4000 o jmp hello (siendo hello una label definida) |
| jeq MEM | 0x21 | si acc = 0, pc = MEM | Salta a la dirección MEM si ACC es 0. |
| jne MEM | 0x22 | si acc != 0, pc = MEM | Salta a la dirección MEM si ACC NO es 0. |
| jgt MEM | 0x23 | si acc > 0, pc = MEM | Salta a la dirección MEM si ACC es mayor que 0. |
| jlt MEM | 0x24 | si acc < 0, pc = MEM | Salta a la dirección MEM si ACC es menor que 0. |
| jfg INT MEM | 0x25 | si flags[INT] == true, pc = MEM | Salta a la dirección MEM si la flag identificada por el número INT esta activa. |

Movimiento
| Instrucción | OpCode | Operación | Explicacón |
|-|-|:-:|-|
| ldr MEM Rx | 0x30 | $mem -> Rx | Carga el contenido de la dirección MEM al registro Rx |
| str Rx MEM | 0x31 | Rx -> $mem | Guarda el contenido del registro Rx a la dirección MEM |
| mov Rx Ry | 0x32 | Rx -> Ry | Copia el contenido de Rx a Ry |
| movi INT Rx | 0x33 | INT -> Rx | Copia el número INT al registro Rx |
| tar RX | 0x34 | ACC -> Rx | Copia el contenido de ACC a Rx |
| tra Rx | 0x35 | Rx -> ACC | Copia el contenido de Rx a ACC |
| inr MEM Rx | 0x36 | memory[$MEM, $MEM+1] -> Rx | Lee la indirección guardada en MEM y la guarda en Rx. Es decir, tenemos en la dirección $4000 el byte 0x31 y en la dirección $4001 el byte 0xb1. Además en la dirección $31b1 tenemos el byte 0xff. Si ejecutamos *inr $4000 R0* estaríamos leyendo las direcciones $4000 y $4001 y buscando en la dirección guardada ahí (en nuestro caso $31b1) para obtener el byte que se guardará en R0. Por lo que al final de esta operación, R0 tendra 0xff. |
| inw Rx MEM | 0x37 | Rx -> memory[$MEM, $MEM+1] | Operación contraria a inr. Guarda el contenido de Rx en la indirección MEM |
| dsk MEM | 0x38 |  | Lee el string guardado en la direccion MEM y lo usa para cargar el fichero binario (rom) con ese nombre. |
| movm MEM0 MEM1 | 0x39 |  | Guarda los dos bytes de MEM0 en MEM1 y MEM1 + 1. |

Llamadas
| Instrucción | OpCode | Operación | Explicacón |
|-|-|:-:|-|
| int INT | 0x40 | pc = 0x0000 + INT | Llama a la interrupción número INT |
| hlt | 0x41 |  | Para la ejecución |
| cll MEM | 0x42 | store state; pc = MEM | Guarda el estado actual de los registros en el stack y salta a la dirección MEM |
| crn | 0x43 | pc = pop() | Recupera el estado de la ejecución antes de la llamada (cll) y continua |
| pmd | 0x44 |  | Habilita el modo protegido |
| ein | 0x45 |  | Habilita las interrupciones |
| din | 0x46 |  | Deshabilita las interrupciones |
| cfg INT | 0x47 | flags[INT] = false | Desactiva la flag identificada por el número INT |

Stack
| Instrucción | OpCode | Operación | Explicacón |
|-|-|:-:|-|
| psa | 0x50 | push(ACC) |  |
| poa | 0x51 | ACC = pop() |  |
| psr Rx | 0x52 | push(Rx) |  |
| por Rx | 0x53 | Rx = pop() |  |

### Ensablado / Desensamblado

Para _"compilar" (ensamblar)_ tu código debes ejecutar tisasm pasando como parámetro el fichero asm que quieras. Puedes desenamblar un binario con el programa tisdiasm y pasando com parámetro la rom.

Por ejemplo, tenemos este código de usuario:

```asm
.data
$5000 "DONE!"

.code $4100
    movi 0x00 R1
    tra R1
:loop
    addi 1
    jfg 0x00 end
    jmp loop

:end
    cfg 0x00
    movm $5000 $0100
	movm $3000 $0102
	int 4
	crn
```

El código lo que hace es suma 1 al ACC hasta que tiene un overflow. Entonces escribe por pantalla el texto "DONE!"

Suponiendo que ese fichero se llama "user.asm", se ensamblaría así:

```
tisasm ./kernal.asm
```

Dejando en el mismo directorio un fichero binario llamado "user.rom".

Para desenamblar se hace con la herramienta tisdiasm

```
tisasm ./user.rom
```

Que saca por pantalla el código asm:

```asm
.data
$5000 "DONE!"

.code $4100
movi 0 R1   		    ;$4100
tra R1   		        ;$4103
addi 1   		        ;$4105
jfg 0 $410e   		    ;$4107
jmp $4105   		    ;$410b
cfg 0   		        ;$410e
movm $5000 $0100   		;$4110
movm $3000 $0102   		;$4115
int 4   		        ;$411a
crn    		            ;$411c
```

Como se puede observar, el código no es exatamente el mismo. Hay unos comentarios en cada línea que indica en qué direccion comienza la instrucción de esa línea. Esto es porque otra diferencia con el código original es que los *labels* no se muestran. En su lugar aparecen direcciones de memoria. Esto es porque los *labels* son sólo una ayuda del ensamblador para el programador, pero al ensamblar, se traducen a las direcciones reales. Por lo que, si se desensambla un código, no se pueden recuperar las *labels*. A cambio, el desensamblador añade las anotaciones con las direcciones para simplificar la lectura del código desensamblado.

## Proceso de arranque
Al iniciar el emulador, lo primero que hace es buscar el binario del kernel, que se debe llamar __kernal.rom__. Hecho esto, lo carga en memoria y comienza a ejecutar las instrucciones a partir de la dirección $0200 (por lo que la sección de código del kernel debe comenzar en esa posición). A partir de este punto se deja completamente el emulador al control del desarrollador del kernel.

## Kernel por defecto
Tienes un pequeño y simple kernel en el fichero *kernel.asm* Asigna algunas interrupciones y te da algunas funcionalidades como la interrupción 4 (strcpy, lee el string de la indirección $0100 y lo deja en la indirección $0102). Además, espera que tu programa se llame *user.rom* y lo ejecuta. Este kernel es usado para el ejemplo de arriba.

Aún así, si quieres puedes mejorar el kernel por defecto o crearte el tuyo propio.