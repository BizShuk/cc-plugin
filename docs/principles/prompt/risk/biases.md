https://www.promptingguide.ai/risks/biases

# Bias Risks Summary

Large Language Models may produce problematic outputs exhibiting biases that reduce performance in downstream tasks.

## Two Main Sources of Bias

### 1. Distribution of Exemplars

In few-shot learning, skewed exemplar distributions affect model predictions.

**Example**: In sentiment classification, flipping the positive/negative exemplar ratio causes the same ambiguous input (like "I feel something") to be classified as different sentiment labels.

### 2. Order of Exemplars

The sequence of exemplars also matters. Placing all positive examples before negative ones may introduce bias, especially when combined with uneven label distribution.

## Best Practices

- Provide balanced number of examples for each label
- Randomize exemplar order
- Conduct extensive experiments to minimize bias