# Rx

[![Discord](https://img.shields.io/discord/428990244952735764.svg?style=flat&logo=discord&colorB=green)](https://discord.gg/M9nxJ3g)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/syncaide/rx.svg?color=e36397&label=release)

### Examples
I will write up a few examples really soon.

### Disclosure
This library is a very similar implementation to the 
[RxGo](https://github.com/ReactiveX/RxGo). The goal of this development is not 
to create a production ready general purpose reactive go library. If you need 
something robust please visit the link and use that library. The major 
difference with that project is that the internals were designed using channels 
and to withstand thread safety of execution. In addition the interface was 
constructed to resemble the way rxjs 6+ works. The purpose of this library is 
to support identical interface for both [gate](https://github.com/syncaide/gate) 
and [chiral](https://github.com/syncaide/chiral). The RxGo contributors 
have done an outstanding work and certainly deserve the credit for it.
