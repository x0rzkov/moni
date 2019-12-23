package main

func init() {
	acl = make(map[string]bool)

	// reject
	acl["localhost"] = false
	acl["google.com"] = false
	acl["github.com"] = false

	// accept
	acl["rustyeddy.com"] = true
	acl["oclowvision.com"] = true
}