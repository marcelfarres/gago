package gago

import (
	"math/rand"
	"sort"
)

// Crossover generates new individuals called "offsprings" are by mixing the
// genomes of two parents.
type Crossover interface {
	Apply(p1 Individual, p2 Individual, rng *rand.Rand) (o1 Individual, o2 Individual)
}

// CrossPoint selects identical random points on each parent's genome and
// exchanges mirroring segments. It generalizes one-point crossover and
// two-point crossover to n-point crossover.
type CrossPoint struct {
	NbPoints int
}

// Apply n-point crossover.
func (cross CrossPoint) Apply(p1 Individual, p2 Individual, rng *rand.Rand) (Individual, Individual) {
	// Choose n random points along the genome
	var (
		points, _ = randomInts(cross.NbPoints, 0, len(p1.Genome), rng)
		nbGenes   = len(p1.Genome)
		o1        = makeIndividual(nbGenes, rng)
		o2        = makeIndividual(nbGenes, rng)
		// Use a switch to know which parent to copy onto each offspring
		s = true
	)
	// Sort the points
	sort.Ints(points)
	// Add the start and end of the genome points
	points = append([]int{0}, points...)
	points = append(points, nbGenes)
	for i := 0; i < len(points)-1; i++ {
		if s {
			copy(o1.Genome[points[i]:points[i+1]], p1.Genome[points[i]:points[i+1]])
			copy(o2.Genome[points[i]:points[i+1]], p2.Genome[points[i]:points[i+1]])
		} else {
			copy(o1.Genome[points[i]:points[i+1]], p2.Genome[points[i]:points[i+1]])
			copy(o2.Genome[points[i]:points[i+1]], p1.Genome[points[i]:points[i+1]])
		}
		// Alternate for the new copying
		s = !true
	}
	return o1, o2
}

// CrossUniformF crossover combines two individuals (the parents) into one
// (the offspring). Each parent's contribution to the Genome is determined by
// the value of a probability p. Each offspring receives a proportion of both of
// it's parents genomes. The new values are located in the hyper-rectangle
// defined between both parent's position in Cartesian space.
type CrossUniformF struct{}

// Apply uniform float crossover.
func (cross CrossUniformF) Apply(p1 Individual, p2 Individual, rng *rand.Rand) (Individual, Individual) {
	var (
		nbGenes = len(p1.Genome)
		o1      = makeIndividual(nbGenes, rng)
		o2      = makeIndividual(nbGenes, rng)
	)
	// For every gene
	for i := 0; i < nbGenes; i++ {
		// Pick a random number between 0 and 1
		var p = rng.Float64()
		o1.Genome[i] = p*p1.Genome[i].(float64) + (1-p)*p2.Genome[i].(float64)
		o2.Genome[i] = (1-p)*p1.Genome[i].(float64) + p*p2.Genome[i].(float64)
	}
	return o1, o2
}

// CrossProportionateF crossover combines any number of individuals. Each of the
// offspring's genes is a random combination of the selected individuals genes.
// Each individual is assigned a weight such that the sum of the weights is
// equal to 1, this is done by normalizing each weight by the sum of the
// generated weights. With this crossover method the CrossSize can be set to any
// positive integer, in other words any number of individuals can be combined to
// generate an offspring. Only works for floating point values.
type CrossProportionateF struct {
	// Should be any integer above or equal to two
	NbParents int
}

// CrossPMX (Partially Mapped Crossover) randomly picks a crossover point. The
// offsprings are generated by copying one of the parents and then copying the
// other parent's values up to the crossover point. Each gene that is replaced
// is permuted with the gene that is copied in the first parent's genome. Two
// offsprings are generated in such a way (because there are two parents). This
// crossover method ensures the offspring's genomes are composed of unique
// genes, which is particularly useful for permutation problems such as the
// Traveling Salesman Problem (TSP).
type CrossPMX struct{}

// Apply partially mixed crossover.
func (c CrossPMX) Apply(p1 Individual, p2 Individual, rng *rand.Rand) (Individual, Individual) {
	var (
		nbGenes = len(p1.Genome)
		o1      = makeIndividual(nbGenes, rng)
		o2      = makeIndividual(nbGenes, rng)
	)
	copy(o1.Genome, p1.Genome)
	copy(o2.Genome, p2.Genome)
	// Choose a random crossover point p such that 0 < p < (nbGenes - 1)
	var (
		p = rng.Intn(nbGenes-2) + 1
		a int
		b int
	)
	// Paste the father's genome up to the crossover point
	for i := 0; i < p; i++ {
		// Find where the second parent's gene is in the first offspring's genome
		a = getIndex(p2.Genome[i], o1.Genome)
		// Swap the genes
		o1.Genome[a], o1.Genome[i] = o1.Genome[i], p2.Genome[i]
		// Find where the first parent's gene is in the second offspring's genome
		b = getIndex(p1.Genome[i], o2.Genome)
		// Swap the genes
		o2.Genome[b], o2.Genome[i] = o2.Genome[i], p1.Genome[i]
	}
	return o1, o2
}
