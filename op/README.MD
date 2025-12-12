## op

### 1. Why You Need `op`

- If you require real-time type-based computations such as:
  - Numeric comparisons
  - Option (enum) lookups
  - Time comparisons
  - Text similarity calculations
  - Conditional merging

If any of the above applies to your use case, we recommend using this package.

### 2. Code Structure

- **`op_bool.go`**: Boolean-type operators and their implementations  
- **`op_has.go`**: Option (enum)-type operators and their implementations  
- **`op_time.go`**: Time-type operators and their implementations  
- **`op_num.go`**: Unary, binary, and ternary numeric-type operators and their implementations  
- **`op_word.go`**: Text-type operators and their implementations, supporting customizable similarity calculation methods:
  - **`default`**: Exact match or substring containment → returns `1`; otherwise → `0`  
  - **`reg`**: Regular expression match → returns `1`; otherwise → `0`  
  - **`levenshtein`**: Levenshtein distance (edit distance) comparing entire sentence against keyword  
  - **`levenshtein-window`**: Sliding-window Levenshtein distance, comparing individual words within a sentence against a keyword