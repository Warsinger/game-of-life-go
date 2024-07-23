# Conway's Game of Life
Conway's Game of Life written in Golang for fun and learning

## Rules
* Any live cell with fewer than two live neighbours dies, as if by underpopulation.
* Any live cell with two or three live neighbours lives on to the next generation.
* Any live cell with more than three live neighbours dies, as if by overpopulation.
* Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

## Controls
* '+' or '-' to adjust game speed
* D enable debug info printout
* R restart the lifecycle

## Command line args
* -cell int
  * Size of each cell in pixels (default 10)
* -debug
  * Show debug info, D to toggle in game
* -height int
  * Board height in cells (default 80)
* -lines
  * Draw grid lines, L to toggle in game
* -speed int
  * Generations per second, min 0 max 60, + or - to adjust in game (default 15)
* -variance float
  * Grid Size (default 0.5)
* -width int
  * Board width in cells (default 80)