# yadab - Yet 'Another' Database

I know.. I know.. Whyyyyyyy, another database? 

Hmmm, for some reason, Im intrigued, often confused, sometimes scared about all the 'engineering' that makes the databases tick. And, its all about getting hands 'dirty', right ;) So, here it is. A labor of 'love'; a rookie relational database from ground up. And, its in `go`. (If I had to do this again, Im unsure if I'll pick go again :thinking: Maybe C# or Java/Kotlin perhaps? lack of 'true' encapsulation in `go` is less than ideal to build something like a database.)

(Maybe I should be doing better with picking a name? doesn't sound as good/natural as yaml :thinking:)

```
  Y   Y      A       DDDD      A       BBBBB
   Y Y      A A      D   D    A A      B    B
    Y      AAAAAA    D   D   AAAAAA    BBBB
    Y     A      A   D   D  A      A   B    B
    Y    A        A  DDDD  A        A  BBBBB

    Welcome! This is `yadab` - yet another database ;)
```

## Features
- Memory Backend: Operates entirely in memory for you know, faster data access. If I invest more time into this, maybe I'll leverage file system.
- Simple Data Types: Supports basic data types such as Text, Int, Float, Double (for now)
- Easy Integration: If this were to become a thing, my hope is that it should be 'easy enough' to embedd or run standalone in various application environments.
- Lightweight: Minimal dependencies and small footprint.
- Extensible: Its rookie atm, but for now, im convinced the structure is good enough to accommodate extension.

## Quick Start

### Pre-requisite
Latest version of `go`

### Installation
Clone the repository:

```
git clone https://github.com/1x-eng/yadab.git
cd yadab
```

Build the project:

```
go build .
```

Run yadab:

```
./yadab
```

### Usage

- Lets create a new table, shall we?

    ```
    CREATE TABLE robots (id INT, name TEXT);
    ```

    should result in something like this - 
    ```
    > ./yadab
      Y   Y      A       DDDD      A       BBBBB  
       Y Y      A A      D   D    A A      B    B 
        Y      AAAAAA    D   D   AAAAAA    BBBB   
        Y     A      A   D   D  A      A   B    B 
        Y    A        A  DDDD  A        A  BBBBB  

        Welcome! This is `yadab` - yet another (relational) database ;)
    # CREATE TABLE robots (id INT, name TEXT);
    ok
    # 
    ```

- And, now, lets seed some data.
    ```
    > ./yadab
      Y   Y      A       DDDD      A       BBBBB  
       Y Y      A A      D   D    A A      B    B 
        Y      AAAAAA    D   D   AAAAAA    BBBB   
        Y     A      A   D   D  A      A   B    B 
        Y    A        A  DDDD  A        A  BBBBB  

        Welcome! This is `yadab` - yet another (relational) database ;)
    # CREATE TABLE robots (id INT, name TEXT);
    ok
    # INSERT INTO robots VALUES (1, 'chatgpt');    
    ok
    # 
    ```

- Then, lets select from our table.
    ```
    SELECT id, name FROM robots;
    ```

    ```
    > ./yadab
      Y   Y      A       DDDD      A       BBBBB  
       Y Y      A A      D   D    A A      B    B 
        Y      AAAAAA    D   D   AAAAAA    BBBB   
        Y     A      A   D   D  A      A   B    B 
        Y    A        A  DDDD  A        A  BBBBB  

        Welcome! This is `yadab` - yet another (relational) database ;)
    # CREATE TABLE robots (id INT, name TEXT);
    ok
    # INSERT INTO robots VALUES (1, 'chatgpt');    
    ok
    # SELECT id, name FROM robots;
    | id | name |
    ====================
    | 1 |  chatgpt | 
    ok
    ```

Well, that's it. Isn't it? What else do you need now :wink:

Maybe, lets insert another robot? why not...

```
> ./yadab
  Y   Y      A       DDDD      A       BBBBB  
   Y Y      A A      D   D    A A      B    B 
    Y      AAAAAA    D   D   AAAAAA    BBBB   
    Y     A      A   D   D  A      A   B    B 
    Y    A        A  DDDD  A        A  BBBBB  

    Welcome! This is `yadab` - yet another (relational) database ;)
# CREATE TABLE robots (id INT, name TEXT);
ok
# INSERT INTO robots VALUES (1, 'chatgpt');    
ok
# SELECT id, name FROM robots;
| id | name |
====================
| 1 |  chatgpt | 
ok
# INSERT INTO robots VALUES (2, 'bard');   
ok
# SELECT id, name FROM robots;
| id | name |
====================
| 1 |  chatgpt | 
| 2 |  bard | 
ok
# 
```

Ok, that's it. I have a life ;)
