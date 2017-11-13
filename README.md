# collections

<!---[![Build Status](https://travis-ci.org/dedis/student_17_collections.svg?branch=develop)](https://travis-ci.org/dedis/student_17_collections)
[![Codecov branch](https://img.shields.io/codecov/c/github/dedis/student_17_collections/develop.svg)](https://codecov.io/gh/dedis/student_17_collections/branch/develop)-->

A `collection` is a Merkle-tree based data structure to securely and verifiably store *key / value* associations on untrusted nodes. The library in this package focuses on ease of use and flexibility, allowing to easily develop applications ranging from simple client-server storage to fully distributed and decentralized ledgers with minimal bootstrapping time.

## Table of contents

- [Overview](#overview)
   * [Basic use example](#basic-use-example)
      + [The scenario](#the-scenario)
      + [The `collection` approach](#the-collection-approach)
   * [Hands on](#hands-on)
      + [Getting started](#getting-started)
      + [Creating a `collection` and a `verifier`](#creating-a-collection-and-a-verifier)
      + [Manipulators](#manipulators)
      + [`Record`s and `Proof`s](#records-and-proofs)

## Overview

### Basic use example

Here we present a simple example problem, and discuss how it can be addressed by using a `collection`, without discussing implementation details, but rather focusing on the high-level features of a `collection` based storage. 

More advanced scenarios will be introduced in later sections of this document. If you prefer to get your hands dirty right away, you can jump to the [Hands on](#hands-on) section of this README.

#### The scenario

The users of an online community want to organize themselves into three groups with different permission levels: `admin`, `readwrite` and `readonly`. Users in the `admin` group can, e.g., read and write documents, and update the permission level of other users (`admin`s included); users with `readwrite` access can, e.g., update documents but not change privileges; and `readonly` users can, e.g., only read documents.

A server is available to the users to store and retrieve information. However, the server is managed by a third party that the users don't trust. Indeed, they will assume that, if it will have the opportunity to do so, the server will maliciously alter the table of permission levels in order, e.g., to gain control of the community.

#### The `collection` approach

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

 * The untrusted server will store one `collection` for each group. Each `collection` will associate no value to each key. Therefore each `collection` will, in practice, represent a set, just like a `map[string]struct{}` can be used to store a set of `string`s.
 * Each user will run a `verifier` for each group, storing only the fixed size *state*.
 * `Update`s to each `collection` will be broadcasted to all `verifier`s, that will check them and update their *state* accordingly.
 * Users will be able to query the server, and use their `verifier`s to verify that it provides honest answers.


![collection](assets/images/collection.gif "Example use scenario")

### Hands on

If you are already familiar with the basic notions behind a `collection`, and just need to get started, this section is what you are looking for. If you need a better understanding of how a `collection` works, you might want to skip this section for now, and dive in [`collection` for dummies](#collection-for-dummies), where you will find all the information you need to bootstrap yourself in the world of `collection`s.

#### Getting started

To install this library, open a terminal `go get` it:

```bash
go get github.com/dedis/student_17_collections
```

Done? Great. Now you just need to import it:

```go
import collections "github.com/dedis/student_17_collections"
```

That's all. You are ready to go. Remember, from now on, we will use `collections` as the name of the package.

#### Creating a `collection` and a `verifier`

The simplest way to get a `collection` and a `verifier` is:

```go
my_collection := collections.EmptyCollection()
my_verifier := collections.EmptyVerifier()
```

This will give you an empty `collection` with no fields (i.e., just a set, as we have seen in [Basic use example](#basic-use-example)), and an empty `verifier`, also with no fields. By empty, here, we mean that the states of `my_collection` and `my_verifier` are the same, and that no record whatsoever can be proved from them.

If you wish to specify one or more value types that you want to be associated to each key, you can provide one or more `field`s to the constructors. You are free to define your own value types (see [Defining custom fields](#defining-custom-fields)) but we provide some that are ready to use.

For example, we could create an empty `collection` (and an empty `verifier`) that to each key associate a 64-bit amount of stake and some raw binary data by

```go
my_collection_with_fields := collections.EmptyCollection(collections.Stake64{}, collections.Data{})
my_verifier_with_fields := collections.EmptyVerifier(collections.Stake64{}, collections.Data{})
```

#### Manipulators

The set of records in a `collection` can be altered using four manipulators: `Add`, `Remove`, `Set` and `SetField`, to add an association, remove one, and set either all or one specific value associated to a key.

In general, you probably want to wrap a set of manipulations on a `collection` in an `Update` object (which also carries proofs: this allows it to be verified and applied also on `collection`s that are not currently storing the records that are being altered, like a `verifier`).

For now, however, let's add some records to our `my_collection_with_fields` and play around with them. Each record, as we said, is an association between a `key` and zero or more values. A `key` is always of type `[]byte` (therefore a `key` can be arbitrarily long!).

The manipulators syntax is pretty straightforward:

```go
my_collection_with_fields.Add([]byte("myfirstrecord"), uint64(42), []byte("myfirstdatavalue")) // Adds an association between key "myfirstrecord" and values 42 and "myfirstdatavalue".
my_collection_with_fields.Add([]byte("anotherrecord"), uint64(42), []byte{}) // Another record with empty data field.

my_collection_with_fields.Remove([]byte("myfirstdatavalue")) // I didn't like it anyway.

my_collection_with_fields.Set([]byte("anotherrecord"), uint64(33), []byte("betterthannothing")) // Note how you need to provide all fields to be set
my_collection_with_fields.SetField([]byte("anotherrecord"), 1, []byte("Make up your mind!")) // Sets only the second field, i.e., the Data one.
```

That's it. Now `my_collection_with_fields` contains one record with key `[]byte("anotherrecord")` and values `(uint64(33), []byte("Make up your mind!"))`.

All manipulators return an `error` if something goes wrong. For example:

```go
err := my_collection_with_fields.Add([]byte("anotherrecord"), uint64(55), []byte("lorem ipsum"))

if err != nil {
	fmt.Println(err)
}
```

Outputs `Key collision.`, since a record with key `[]byte("anotherrecord")` already exists in the `collection`. 

An `error` that all manipulators can return is `Applying update to unknown subtree. Proof needed.`. This happens when you try to manipulate a part of a `collection` that is not locally stored. Remember that, while the allowed set of records a `collection` can store is uniquely determined by its *state*, a `collection` can store an arbitrary subset of them. In particular, as we said, `verifier`s store no association locally, so if you try to

```go
err = my_verifier_with_fields.Add([]byte("myrecord"), uint64(65), []byte("Who cares, this will not work anyway."))

if err != nil {
	fmt.Println(err)
}
```

prints out `Applying update to unknown subtree. Proof needed.`. To manipulate records that are not locally stored by a `collection`, you first need to verify a `Proof` for those records. But all in due time, we will get to that later.

#### `Record`s and `Proof`s

Now that we know how to add, remove and update records, we can discuss what makes a `collection` such a useful instrument: `Proof`s.

##### `Get()`ting a `Record`

Let's first start with something easy. As we said, a `collection` is a *key / value* store. This means that in principle we could use a `collection` as we would use an (unnecessarily slow) Go `map`.

As in any *key / value* store, data can be retrieved from a `collection` by providing a `key`. Following from the `my_collection_with_fields` example, we can try to retrieve two `Record`s, one existing and one non-existing, by:

```Go
another, anothererr := my_collection_with_fields.Get([]byte("anotherrecord")).Record()
nonexisting, nonexistingerr := my_collection_with_fields.Get([]byte("nonexisting")).Record()
```

Now, **attention please!** Here one could expect `anothererr` to be `Nil` and `nonexistingerr` to be some `Error` that says something like: `Key not found`. However, *this is not what those errors are for!* When fetching a `Record`, not finding any association corresponding to a `key` here is not an error: indeed, the retrieval process was completed successfully and no association was found. So both of the above calls will return a `Record` object and no error. 

However, we *can* use `Match()` to find out if an association exists with keys `anotherrecord` and `nonexisting` respectively by

```Go
fmt.Println(another.Match()) // true
fmt.Println(nonexisting.Match()) // false
```

What call to `Get()` can return an `Error`, then? Simple: an `Error` is returned when a `collection` *does not know* if an association exists with the `key` provided because it is not storing it. For example a `verifier`, which does not permanently store any association, when queried as follows:

```Go
another, anothererr = my_verifier_with_fields.Get([]byte("anotherrecord")).Record()
fmt.Println(anothererr)
```

prints out: `Record lies in an unknown subtree.`

Along with `Match()`, a `Record` object offers two useful getters, namely:

```Go
fmt.Println(another.Key()) // Byte representation of "anotherrecord" (just in case you forgot what you asked for :P)

anothervalues, err := another.Values() // Don't worry, err will be non-Nil only if you ask for Values() and Match() is false.
fmt.Println(anothervalues[0].(uint64)) // 33
fmt.Println(anothervalues[1].([]byte)) // Byte representation of "Make up your mind!"
```

Please notice the type assertions! Since each field can be of a different type, `Values()` returns a slice of `interface{}` that need type assertion to be used.

Let us now stop for a second and underline what a `Record` is. A `Record` is just a structure that wraps the result of a query on a `collection`. It just says: "there is an association in this `collection` with key `anotherrecord` and values `33` and `Make up your mind!`", or: "there is no association in this `collection` with key `nonexisting`". You can **trust** that information only because, well, it was generated locally, by a computer that runs your software. However, what if, e.g., someone sent you a `Record` object over the Internet? You couldn't trust what it says more than any other message.

That is why we need `Proof`s.

##### `Get()`ting a `Proof`

So far, we have encountered `verifier`s quite many times, and you are probably wondering what purpose do they serve: they cannot be manipulated, you cannot `Get()``Record`s out of them... apparently, you can only create one and have it sit there for no reason!












