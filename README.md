# Intrusive containers

This library provides multiple implementations of various container data structures written in Go.  
All data structures here are **intrusive** meaning that metadata for managing a data structure is not hidden by implementation but is explicitly located inside user's data types.  
This library is inspired by [Boost.Intrusive](https://www.boost.org/doc/libs/1_84_0/doc/html/intrusive.html) library for C++ programming language.  

## Containers

There are three broad categories of containers represented here:  
1. Hash based - in general operations are performed with amortized constant complexity
    1. [planned] `Map` - a mapping of key to values
    1. [planned] `Set` - a representation of set of objects
1. Tree based - in general operations are performed with logarithmic complexity
    1. `RbTree` - self-balancing binary search tree, it's very similar to a concept of a set
    1. [planned] `MapTree` - self-balancing binary search tree that holds mapping of key to values
1. Lists - complexity varies based on type of operation
    1. `SList` - singly-linked list
    1. `DList` - doubly-linked list, queues could be implemented on top of it

## Pros & Cons

It is strongly suggested to refer to [documentation](https://www.boost.org/doc/libs/1_84_0/doc/html/intrusive/intrusive_vs_nontrusive.html) of Boost.Intrusive for explanation of benefits and shortcomings of using intrusive containers but the main selling point is that one object can contain metadata for multiple containers and as so can be placed in multiple containers at once witch facilitates memory reuse and cache utilization.  

## Design

Design goals for this library are:  
* One object should be able to be inserted in as many containers as amount of embedded metadata it has
* Be functionally complete based on existing containers in [Go](https://pkg.go.dev/container/list) and [C++](https://en.cppreference.com/w/cpp/container)
* Container interface should be similar across all implemented containers
* Some containers in Boost offer optimization of size field. This library do not implement such optimization and as a result require head of the structure to be present for most modification actions
* Container should be tested with fuzz-testing with 97-100 % line coverage
* Container should be tested with unit tests with 80+ % line coverage

## Attribution

This project was developed with the assistance of a large language model (LLM) DeepSeek-V3 noreply@deepseek.com.  

Design of implemented data structures is well known and there is a lot of material that was probably used as training data for LLMs. So there is a reason to believe that generated code would be of high quality. Everything should be tested to proof the quality of implementation. The end-goal of a project is to have intrusive containers in Go and not for authors to learn the skill of implementing data structures from scratch.  

In the light of what written above the usage of LLMs is a logical step towards project completion.
