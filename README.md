# collections

[![Build Status](https://travis-ci.org/dedis/student_17_collections.svg?branch=develop)](https://travis-ci.org/dedis/student_17_collections)
[![Codecov branch](https://img.shields.io/codecov/c/github/dedis/student_17_collections/develop.svg)](https://codecov.io/gh/dedis/student_17_collections/branch/develop)

A `collection` is a Merkle-tree based data structure to securely and verifiably store *key / value* associations on untrusted nodes. The library in this package focuses on ease of use and flexibility, allowing to easily develop applications ranging from simple client-server storage to fully distributed and decentralized ledgers with minimal bootstrapping time.

## Table of contents

- [Overview](#overview)
   * [Basic use example](#basic-use-example)
      + [The scenario](#the-scenario)
      + [The `collection` approach](#the-collection-approach)

# Overview

## Basic use example

Here we present a simple example problem, and discuss how it can be addressed by using a `collection`, without discussing implementation details, but rather focusing on the high-level features of a `collection` based storage. 

More advanced scenarios will be introduced in later sections of this document. If you prefer to get your hands dirty right away, you can jump to the [Hands on](#hands-on) section of this README.

### The scenario

The users of an online community want to organize themselves into three groups with different permission levels: `admin`, `readwrite` and `readonly`. Users in the `admin` group can, e.g., read and write documents, and update the permission level of other users (`admin`s included); users with `readwrite` access can, e.g., update documents but not change privileges; and `readonly` users can, e.g., only read documents.

A server is available to the users to store and retrieve information. However, the server is managed by a third party that the users don't trust. Indeed, they will assume that, if it will have the opportunity to do so, the server will maliciously alter the table of permission levels in order, e.g., to gain control of the community.

### The `collection` approach

As we will see in the next sections, organizing data in a `collection` allows its storage to be outsourced to an untrusted third party. In short:

 - Every `collection` object stores, at least, an *O(1)* size *state* (namely, the label of the root of a Merkle tree, see later).
 - The *state* of a `collection` uniquely determines **all** the *key / value* associations that the `collection` *can* store. However, a `collection` can actually store an arbitrary subset of them. In particular, we will call `verifier` a `collection` that stores no association (note how a `verifier` is therefore *O(1)* in size).
 - A `collection` storing an association can produce a `Proof` object for that association. A `Proof` proves that some *association* is among those allowed by some *state*. Provided with a `Proof` that starts with its *state*, a `collection` can verify it and, if it is valid, use it to add the corresponding associations to the ones it is storing.
- A `collection` can alter the associations it is storing by adding, removing and changing associations. This, however, changes the *state* of the `collection`.
- One or more `Proof`s can be wrapped in an `Update` object. Provided with an `Update` object, a `collection` can:
   * Verify all the proofs contained in it and use them, if necessary, to extend its storage with the associations it wasn't storing.
   * Verify that some update is applicable on the associations contained in the `Proof`. What constitutes a valid `Update` clearly depend on the context where `collection`s are used. Indeed, a ledger can be defined as a set of allowed `Update`s.
   * Apply the `Update`, thus producing a new *state* for the `collection`.
   * Drop all the associations it is not willing to store. Remember that *as long as it keeps track of the state, the changes applied to the associations will be permanent!*.

Our users can therefore make use of the untrusted server as follows:

 * The untrusted server will store one `collection` for each group. 
 * Each user will run a `verifier` for each group, storing only the fixed size *state*.
 * `Update`s to each `collection` will be broadcasted to all `verifier`s, that will check them and update their state accordingly.
 * Users will be able to query the server, and use their `verifier`s to verify that it provides honest answers.


![collection](assets/images/collection.gif "Example use scenario")