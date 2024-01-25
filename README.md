# GRAB

A CLI program that searches for patterns in text.

## Overview

The CLI searches patterns passed as positional arguments in text passed via piping and in attached files.

It uses _Boyer-Moore-Horspool_ algorithm for searching pattern in a string under the hood.

## Installing

Install with `go install`

`go install github.com/danblok/grab/cmd/grab@latest`

Check out available releases [Releases](https://github.com/danblok/grab/releases)

The latest release here [latest](https://github.com/danblok/grab/releases/tag/v0.1.2)

## How to use

### Piping

You can pipe text and use the program like this

![image](https://github.com/danblok/grab/assets/91749788/a2d76989-aa83-4e32-88e3-7a3214a063b5)

### Attaching files

Use `--files` flag to attach files to search patterns there

![image](https://github.com/danblok/grab/assets/91749788/4987bd1c-8871-49b2-9c26-cdd26d00e083)

### Both

Also you can search for patterns using piping and attaching files simultaneously.

![image](https://github.com/danblok/grab/assets/91749788/7fdcc078-be53-49d9-a656-654e6e1487a1)

### Quite mode

Output in quite mode similar to `grep` command.

![image](https://github.com/danblok/grab/assets/91749788/d52ef276-6dcd-4b56-b275-1162d2c19865)

### Non-human mode

Prints line number as the first number, then prints pairs of the start and the end positions of found patterns.

![image](https://github.com/danblok/grab/assets/91749788/40a4dcd3-739f-41bc-8f69-a69028180a57)
