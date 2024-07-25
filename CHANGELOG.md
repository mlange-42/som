# Changelog

## [[unpublished]](https://github.com/mlange-42/som/compare/v0.1.0...main)

### Features

* Adds label propagation for semi-supervised learning (#43)
* Adds switch `--ignore` to all prediction-related commands, for ignoring layers (#43)
* Adds weight decay for better regularization during training (#44)
* Can handle missing/no-data class labels (#45)
* Adds CLI command `fill` for filling missing data (#49)
* Adds CLI command `predict` for predicting complete layers (#50)
* Adds `codes` plots showing a 2D array of small line, pie, rose or image plot, one for each node (#52, #53, #53, #54, #55, #57)
* Adds CLI flag `--sample` to plot only a sub-sample of the input data (#56)

### Documentation

* Examples are included in release downloads (#51)
* Adds hand-written digits example (#55)
* Adds world countries population pyramids example (#57)

### Bugfixes

* Update BMU weights before the weights of neighboring nodes (#42)

### Other

* All plots draw the legend on top of / after the actual plot (#48)

## [[v0.1.0]](https://github.com/mlange-42/som/commits/v0.1.0/)

Initial release of SOM, the Self-organizing Maps library and CLI tool for [Go](https://go.dev).
