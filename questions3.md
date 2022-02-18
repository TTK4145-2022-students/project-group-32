Exercise 3 - Project brainstorm
-------------------------------

## other stuff
We should have a time delay before the elevator tries to take an order. So it can be stoped by a better offer from antother elevator.


## Implementation
What is a module made of?
> Files. We can allso have modules (a folder) made of submodules (folder/file)

How will the modules interact?
> Functions, message passing and data.json.

## Data from outside world
Are the network messages different
> No we're planning on only sending one type. But the modules should be generic

Should the program do different things when receiving different kinds of messages?
> We only have one msg type

Should the different kinds of data feed into the same "top level" module?
> No just write/update to file (through some acceptance)

## The contents
Briefly, and as "seen from the outside": What modules do you need?
> Main (run go routines)
Order processor / control unit:

	- Cab state
	- Order state (state machine and data)
	- Offer algorithm
	- Priotity algorithm

Network

	- UDP_send
	- UDP_recieve

Filesystem

	- Encoding/write
	- Decoding/read

Controll string
Phoenix

Hardware

	- Input panel
	- Lights
	- Cab
		○ Sensors (read)
		○ Motor
		○ Door

In order for a module to perform its task, it depends on information. What are the inputs, outputs, and state of the modules? 
> Orders table, elevator state

Are you playing to strengths of language
> Using go routines and mutex


