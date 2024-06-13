package learn

import "math"

// neuralNet contains all of the information
// that is needed to run a neural network
type neuralNet struct {
	config  neuralNetConfig
	wHidden [][]float64
	bHidden []float64
	wOutput [][]float64
	bOutput []float64
}

// neuralNetConfig contains the configuration
// for a neural network.
// It contains the number of input neurons,
// the number of hidden neurons (int), the number of output neurons (int),
// the number of epochs (int), and the learning rate.
type neuralNetConfig struct {
	// inputNeurons is the number of input neurons
	// in the neural network.
	inputNeurons int
	// hiddenNeurons is the number of hidden neurons
	// in the neural network.
	hiddenNeurons int
	// outputNeurons is the number of output neurons
	// in the neural network.
	outputNeurons int
	// numEpochs is the number of epochs to train the neural network.
	// The neural network will train for this number of epochs.
	numEpochs int
	// learningRate is the learning rate to use for training the neural network.
	learningRate float64
}

// newNeuralNet creates a new neural network
// with the given configuration.
func newNeuralNet(config neuralNetConfig) *neuralNet {
	return &neuralNet{
		config: config,
	}
}

// sigmoid is a sigmoid function implementing the sigmoid activation function.
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// sigmoidPrime is the derivative of the sigmoid function.
// It is used to calculate the gradient of the sigmoid function during
// backpropagation.
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}
