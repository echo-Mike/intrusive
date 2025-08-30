# Intrusive containers

This library provides multiple implementations of various container data structures written in Go.  
All data structures here are **intrusive** meaning that metadata for managing a data structure is not hidden by implementation but is explicitly located inside user's data types.  
This library is inspired by [Boost.Intrusive](https://www.boost.org/doc/libs/1_84_0/doc/html/intrusive.html) library for C++ programming language.  

## Containers

There are three broad categories of containers represented here:  
1. Hash based - in general operations are performed with amortized constant complexity  
1. Tree based - in general operations are performed with logarithmic complexity  
1. Lists - complexity varies based on type of operation  

## Pros & Cons

It is strongly suggested to refer to [documentation](https://www.boost.org/doc/libs/1_84_0/doc/html/intrusive/intrusive_vs_nontrusive.html) of Boost.Intrusive for explanation of benefits and shortcomings of using intrusive containers but the main selling point is that one object can contain metadata for multiple containers and as so can be placed in multiple containers at once.  

## Design

Design goals for this library are:  
* One object should be able to be inserted in as many containers as amount of embedded metadata it has
* Be functionally complete based on existing containers in [Go](https://pkg.go.dev/container/list) and [C++](https://en.cppreference.com/w/cpp/container)
* Container interface should be similar across all implemented containers

## Attribution

This project was developed with the assistance of a large language model (LLM) DeepSeek-V3 noreply@deepseek.com.
