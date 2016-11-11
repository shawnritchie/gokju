GoKju
------
In the process of learning GoLang and other programming techniques I am trying to develop a toolkit to be used in a CQRS/ES applications.

Please keep in mind that at this point the project is being used as a learning tool, so keep the comments constructive :D

Packaging Lessons Learnt
-------------------------
- Multiple files can co-exist under the same package which would simulate having everything under one file
    - [ ] check if this is good practise in the go community
- Watch out for circular dependencies (Still need to investigate how one can avoid this)

Interfaces Lessons Learnt
--------------------------
- Quite a powerful construct in the language, currently trying to make heavy us of the feature to make the code easily extendable

Reflections Lessons Learnt
---------------------------
- Interfaces are represented by a tuple (reflect.Type, reflect.Value)
- Interfaces can be implemented either by a pointer to a struct or by a struct, either one or the other is implementing the interface not both

TODO::
------
- [ ] Implement logging
- [ ] Implement unit test for the current functionality
- [ ] Implement a asynchronous router, making use of buffered channels, guranteeing order but synchronous execution of the messages
- [ ] Implement a circuit breaker which fails all message delivery atempts in the channel bugger is full
- [ ] Implement life cycle manager able to pause, shutdown and signal
- [ ] Implement clean shutdown functionality making sure all synchronous messages are delivered and that all items on the mailbox are consumed
- [ ] Implement the concept of unit of work allowing multiple transactions to be managed as one
