# ASCII-Art-Web

## Description

ASCII-Art-Web is a Go-based web application that allows users to generate ASCII art directly from their browser. Instead of printing output to the terminal, the application provides a graphical interface where users can enter text, choose a banner style, and view the generated ASCII art instantly.

The application supports multiple banner styles:

* standard
* shadow
* thinkertoy

---

## Author

Arogundade Adewale Muideen

---

## Usage

### 1. Start the server

```bash
go run .
```

### 2. Open in browser

```
http://localhost:8080
```

### 3. Generate ASCII Art

* Enter text in the input field
* Select a banner style
* Click submit
* View the ASCII art result on the page

---

## Implementation Details

### Server

* Built using Go’s `net/http` package
* Runs on port `:8080`
* Handles routes:

  * `/` → renders the main page
  * `/ascii-art` → processes form submission

---

### Banner Processing

* Banner files are read from disk
* Each character is mapped to its ASCII representation
* Stored in a map structure:

```
map[rune][]string
```

---

### ASCII Art Generation

* Input string is processed character by character
* Each character is looked up in the banner map
* ASCII art is constructed row-by-row:

  * First row of all characters
  * Then second row, and so on

This preserves alignment and structure.

---

### Display (Rendering)

* Uses Go’s `html/template` package
* Data is passed from the server to the template
* ASCII output is rendered inside a `<pre>` tag to preserve spacing

---

## Notes

* Only supported ASCII characters are processed
* Invalid input or missing banners are handled with error responses
* Output formatting depends on the integrity of banner files

---
