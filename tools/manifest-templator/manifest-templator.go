/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type templateData struct {
	Namespace          string
	ContainerTag       string
	ContainerPrefix    string
	ImagePullPolicy    string
	GeneratedManifests map[string]string
}

func main() {
	namespace := flag.String("namespace", "", "")
	containerPrefix := flag.String("container-prefix", "", "")
	containerTag := flag.String("container-tag", "", "")
	imagePullPolicy := flag.String("image-pull-policy", "IfNotPresent", "")
	genDir := flag.String("generated-manifests-dir", "", "")
	inputFile := flag.String("input-file", "", "")
	flag.Parse()

	data := templateData{
		Namespace:          *namespace,
		ContainerTag:       *containerTag,
		ContainerPrefix:    *containerPrefix,
		ImagePullPolicy:    *imagePullPolicy,
		GeneratedManifests: make(map[string]string),
	}

	if *genDir != "" {
		manifests, err := ioutil.ReadDir(*genDir)
		if err != nil {
			panic(err)
		}
		for _, manifest := range manifests {
			b, err := ioutil.ReadFile(filepath.Join(*genDir, manifest.Name()))
			if err != nil {
				panic(err)
			}
			data.GeneratedManifests[manifest.Name()] = string(b)
		}
	}

	tmpl := template.Must(template.ParseFiles(*inputFile))
	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
