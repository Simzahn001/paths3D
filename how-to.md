## 1. Create a grid

To start of, we have to create a grid. There are several ways of creating a grid:

- **Create an empty grid:**

  Specify the size of the grid, how many cells the grid should be wide.
  In this case `x=10` and `y=10`.

    ```go
    var grid = NewGrid(10,10,1,1)
    ```

- **Create a grid from a string array, which gives the gird a layout:**

  Each letter stands for a cell. These letters will be saved in each cell's `rune` parameter.
  You don't have to specify the size, it will be read from the input.

    ```go
    var layout = []string{
        "aabqponm",
        "abcryxxl",
        "accszzxk",
        "acctuvwj",
        "abdefghi",
  }
    var grid = NewGridFromStringArray(layout, 1, 1)
    ```    

## 2. Edit parameters

You can edit each cell's parameters like: walkable, height, cost and rune.

These values can be edited for one cell alone, or for multiple cells:

- **Edit a parameter for one specific cell:**

  Get the cell by the position `x=2` and `y=3` and set its rune to "t".
  This is done just by using the basic setter of the cell struct.
  ```go
  grid.Get(2, 3).Rune = 't'
  ```
  
-------------------------------------

- **Edit multiple cells**

  You can change the parameter of all cells with the same cell-rune at once.

    <br>

  This sets: all `x`-cells to be walkable:
  ```go  
  grid.SetWalkable('x', true)
  ```
  <br>

  The same can be done with other parameters like `cost` :
  ```go
  grid.SetCost('g', 10)
  ```
  <br>

  > **AddHeightMap()**
  > 
  > Especially for grids, which represent something like a terrain or landscape, `AddHeightMap()`
  can be useful.
  >
  > If we take the example from above (a grid created from a string array), this string array may
  stand for a landscape, where `a` is the lowest and `z` is the highest point; ==> a hill
  > ```go
  > var allLetters //an array containing all 26 letters of the latin alphabet
  >
  > heightMap := make(map[rune]int)
  > 
  > for i, letter := range allLetter {
  > eightMap[letter] = i
  > }
  >  
  > grid.AddHeightMap(heightMap)
  > ```

## 3. Pathfinding

The last step of creating the path itself is pretty straight-forward. Simply call `GetPathFromCells` with the
cell to start from and to cell to go to:

```go
grid.GetPathFromCells(grid.Get(1,1), grid.Get(3,5), false, false)
```

<br>

The tow additional parameters do set some limitations to the pathfinding:

1. **Diagonals**

    If this parameter is set to true, diagonal movements will be allowed.

    <br>

    In this example, a path from the bottom left corner to the top right corner should be found.
    This should represent a 5x5 grid, with all walkable & the same hight cells (dots `.`).
    The path is shown using the hashtag (`#`).
    ```
    //  true       false
      +-----+     +-----+
      |....#|     |#####|
      |...#.|     |#....|
      |..#..|     |#....|
      |.#...|     |#....|
      |#....|     |#....|
      +-----+     +-----+
    ```
    <br>

2. **Wall blocks diagonals**
    
   If this parameter is set to true, diagonal movements will be able "trough" walls. If diagonals is disabled,
   this setting doesn't have any impact.

    <br>

    Again, a path from the bottom left to the top right is searched for.
    The dots represent an "empty" cell, while hashtags do represent a wall (not walkable cell).
    ```
    //  true       false
      +-----+     +-----+
      |.####|     |....#|
      |#x...|     |.x.#.|
      |#x...|     |.x#..|
      |#.xxx|     |.#xxx|
      |#....|     |#....|
      +-----+     +-----+
    ```