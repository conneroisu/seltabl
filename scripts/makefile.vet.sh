#!/bin/bash

go vet ./...
sqlc vet
cd ..
