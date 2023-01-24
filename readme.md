
# paths3D

## What is paths3D?

paths3D is a pathfinding library written in Golang. Its main feature is simple best-first and shortest-cost path finding.


## Fork goals
This is a fork form [SolarLune's paths](https://github.com/SolarLune/paths) path finding library.  The fork aims 
at adding a third dimension to the pathfinding system: height. Each cell will have a specific height level.
By default, the pathfinding only can go 1 height unit up per cell. And can drop down infinite Blocks. 
Configuration is planned ;) 

## Fork Status
| Feature                                                                               |  State   |
|---------------------------------------------------------------------------------------|:--------:|
| **Height Pathfinding**<br>Add pathfinding with heiht levels                           | Finished |
| **Step Heights**<br>Make the height that can be stepped up configurable               | Finished |
| **Drop Heights**<br>Make the height that can be dropped down at once configurable     | Finished |


## Why did you create paths3D?

I was doing day 12 of the [AoC](https://adventofcode.com/2022/day/12) and was wondering if there is a simple 
and straight-forward pathfinding library for this. I couldn't find a lib for this use-case, so I decided to create one. 
Because I didn't want to start from scratch, so I decided to fork the existing paths lib and extend it with a 
height level.

## How do I install it?

Just go get it:

`go get github.com/Simzahn001/paths3D`

## How do I use it?

Basically paths3D is based on a grid with cells. Each cell does have a position, walkability and height. Further 
you can specify a cost; lower price cells will be preferred during the path finding process.

To find a detailed and explained example to start, please take a look at the [how-to.md](https://github.com/Simzahn001/paths3D/blob/master/how-to.md).

Here is a whole example of the lib:

```go

import "github.com/Simzahn001/paths3D"

func Init() {
	
    //create a grid with a string array
    layout := []string{
        "xxxxxxxxxx",
        "x        x",
        "x xxxxxx x",
        "x xg   x x",
        "x xgxx x x",
        "x gggx x x",
        "x xxxx   x",
        "x  xgg x x",
        "xg ggx x x",
        "xxxxxxxxxx",
    }
    grid := paths.NewGridFromStringArrays(layout)

    // After creating the Grid, you can edit it using the Grid's functions. In this case, we make the
	// cells to "walls". Note that here, we're using 'x' to get Cells that have the rune for the lowercase
	//x character 'x', not the string "x".
    grid.SetWalkable('x', false)

    // You can also loop through them by using the `GetCells` functions thusly...
    for _, goop := range grid.GetCellsByRune('g') {
        goop.Cost = 5
    }

    // This gets a new Path from the Cell occupied by a starting position [24, 21], to another [99, 78]. The next two Parameters
	// specify the maximum height that can be stepped up at once and the maximum height that can be dropped down at once.
	// The last two parameters are settings for diagonal movement between two cells.
    grid := GameMap.GetPathFromCell(GameMap.Get(1, 1), GameMap.Get(6, 3), 1, 5 false, false)

    // After that, you can use Path.Current() and Path.Next() to get the current and next Cells on the Path. When you determine that 
    // the pathfinding agent has reached that Cell, you can kick the Path forward with path.Advance().
	
}
```

## Dependencies?

For the actual package, there are no external dependencies.


## SolarLune's path docs:
[GoDocs](https://pkg.go.dev/github.com/SolarLune/paths?tab=doc)
