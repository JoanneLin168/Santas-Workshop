# Santa's Workshop Simulation #

![Santa's Workshop](SantasWorkshop.png)

This is a project to simulate Santa and his elves preparing the presents for children for Christmas.

This project is written in Golang and uses the [Ebiten library](https://github.com/hajimehoshi/ebiten) for the visualisation,
and Chris Gora's [semaphore](https://github.com/ChrisGora/semaphore) library.

The following were used in this project:
- Multithreading (to simulate the elves)
- Remote Procedure Calls (RPC) (to send the list of children to Santa's Workshop)
- Mutex locks (to prevent any race conditions)
- Semaphores (to simulate the maximum number of elves in the storage room)

Watch this to see it working!
<iframe width="560" height="315" src="https://youtu.be/RYMBpJ0iSx8" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe> 

## How to Run ##
Currently there is no server set up, so you will need to run the server locally.
To do this, navigate to the `server` folder in a terminal, then type `go run workshop.go`.
Next, open up another terminal and navigate to the project directory, then type `go run main.go`.

You can also run tests by navigating to the `test` folder and type `go test -v`.
You may want to consider running a specific test, which you can do by typing `go test -v -run=` followed by one of the following tests:
- TestResults
- TestDistances
- TestInput
- TestTime
