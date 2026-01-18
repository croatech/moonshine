# Agents Guidelines

## Code Style Rules

### Comments

**NEVER write comments in code.**

Code should be self-documenting through:
- Clear variable and function names
- Well-structured code
- Proper architecture and separation of concerns

If code needs explanation, it should be refactored to be more readable instead of adding comments.

Exceptions:
- Public API documentation (godoc comments for exported functions/types)
- License headers
- Generated code markers

Don't put debug comments like this [LocationService] Failed to commit transaction etc







