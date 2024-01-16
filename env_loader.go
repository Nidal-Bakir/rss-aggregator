package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func readAndSetEnv(envPath string) {
	file, err := os.OpenFile(envPath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("can not read .env file")
	}

	v := bufio.NewScanner(file)

	for v.Scan() {
		key, val, ok := strings.Cut(v.Text(), "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		os.Setenv(key, val)
	}

	_, isExist := os.LookupEnv("PORT")

	if !isExist {
		log.Fatal("The PORT env not exist in the .env file, A porting...")
	}
}
