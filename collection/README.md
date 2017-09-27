# collection

## Overview

Package `collection` is a Golang implementation of a Merkle-tree based (optionally weighted) collection.

A `collection` is an authenticated data structure that allows one or more `verifier`s to outsource the storage of a set of identifiers to one (or more) untrusted servers.

### Example use scenario

A distributed service wants to organize its users into three groups with different permission levels: `admin`, `readwrite` and `readonly`. Users in the `admin` group can, e.g., read and write documents, and update the permission level of other users (`admin`s included), users with `readwrite` access can, e.g., update documents but not change privileges, and `readonly` users can, e.g., only read documents.

In order to do this, the three sets of user identifiers (e.g., their public keys) are stored on an untrusted server, that users can query and push updates to. Each users should, however, have the possibility to check if the server is honestly storing the sets, without having to store a complete copy of them.

This can be achieved using `collection`s: 

 * The untrusted server will store one `collection` for each group. 
 * Each user will run a `verifier` with a fixed size state.
 * Updates to the `collection` will be broadcasted to each `verifier`, that will check them and update its state accordingly.
 * Users will be able to query the server, and use their `verifier`s to verify that it provides honest answers.

 
![collection](https://raw.githubusercontent.com/dedis/student_17_collections/develop/collection/assets/images/collection.gif "Example use scenario")
