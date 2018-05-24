#!/bin/bash
go build && cf install-plugin -f cf-all-services && cf all-services
