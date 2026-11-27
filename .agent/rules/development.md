---
trigger: always_on
---

# **Important!** Rules of Development (Code book)

## 1. Core Principles
- **KISS - Keep it simple, stupid** - If something can be done more simply, do it more simply
- **YAGNI - You aren't gonna need it** - Don't write code "just in case it might be useful later", don't try to predict all future uses of your code, just do what's needed now
- **DRY - Don't repeat yourself** - Don't copy longer code fragments. If you need similar code to what exists elsewhere, make the existing place more universal

These 3 principles are equally important; we follow all three simultaneously.

## 2. Variable/Function/Class Names
Names in the code should be long and self-explanatory. Today's IDEs suggest variable names as we type, so using longer names doesn't require more keyboard typing. Long names have many advantages. They often replace comments and specifications. Well-named code is immediately readable for someone new looking at a snippet of the code. AI tools like Copilot will also better suggest subsequent lines of code when everything is named with full names, as AI can more easily infer context (just like humans). We'll also avoid accidental conflicts of functions with the same names or human mistakes caused by similar names.

**Good examples of naming:**
- addTwoArrays()
- copyMarketplaceAccountSettings()
- activeAccountsUserIds

**Bad examples of naming:**
- add(), ATA(), addArr()
- copy(), CMAS(), copyStg()
- active, ids, users, AAUI

## 3. File Lengths (line count)
- Function should fit within 100 lines
- File should contain no more than 1000 lines
- File should contain at least 20 lines of actual logic (not counting lines with braces or empty lines)

These are soft limits; in justified cases, longer/shorter code can be used. These are suggested guidelines. Of course, there are already many longer files that were much smaller when created but grew over time. These are the main candidates for mini-refactors and splitting.

## 4. Constants
Magic ints and magic strings are evil. We don't want to use them; we always create constants. This makes it easier to search where a constant is used and significantly improves code readability.


------------------

# Code Writing and Code Review Best Practices

Every piece of code should be readable and clear. It should be understandable not only in the context of the entire file but also within smaller code fragments.

To achieve this, eliminate any aspect that requires additional analysis from the developer. Every extra element for analysis contributes to decreased code readability and understanding.

Adhering to the following rules will help address these issues.

## 1. Code Organization

Before adding any code, consider its appropriate placement.

*   **Logical Blocks**: Related methods should be located next to each other.
*   **Cause-and-Effect Chain**: If one method calls another, the latter should, if possible, be located immediately after the first.
*   **Divide and Conquer**: For complex logic, divide large functions into smaller, more specialized functions. Naming these smaller methods appropriately aids in understanding the overall algorithm.
*   **Create on Demand**: Variables should be created (and their values determined) just before their first use. This avoids scrolling to find variable definitions and prevents unnecessary value assignments if the function exits early.

## 2. Variable Naming

Good naming of variables and methods increases code readability and often eliminates the need for additional comments.

*   **Clarity**: Naming should clearly indicate the variable's content or the method's action.
*   **Consistency**: Use the same naming convention for a given type of data, structure, or interface element throughout the codebase. Inconsistent naming (e.g., `field` / `fieldKey` / `fieldId` / `fieldName`) hinders code reading and requires additional analysis.
*   **Verbs as Prefixes for Boolean Variables/Methods**: Names of variables and methods that store/return booleans should start with a verb indicating their state or action, e.g., `is`, `was`, `has`, `should`.

## 3. Less Code is Better

Minimizing code reduces the amount of text to analyze and potential error points. Smaller functions are also easier to test. This should be done without sacrificing readability or performance.

*   **DRY (Don't Repeat Yourself)**: Extract common parts from multiple code fragments.
*   **Code Minimization**: If the same functionality can be achieved with less code without reducing readability, do so.
*   **Removing Redundant Code**: Eliminate unnecessary `if` statements, function parameters, or variables.
*   **Cleaning Code Before Commit**: Before creating a pull request, review your code and remove temporary remnants (e.g., variables named `$tmp`, `var_dumps`, debuggers, unnecessary comments).

## 4. Method Parameters

Parameters should be clear and easy to use. A high number of parameters often suggests the method should be split. Each additional parameter requires additional analysis.

*   **Minimize Parameter Count**: Eliminate method parameters if their values can be determined within the method without negatively impacting performance.
*   **Avoid Boolean Parameters**: Boolean parameters reduce readability at the call site and can indicate that a function handles two independent processes. Consider splitting such functions.

## 5. Avoiding Nested Conditions

Excessive indentation reduces code readability. Avoid deeply nested `if` statements ("pyramid of doom") by using patterns like "Return Early" or by extracting nested code into separate functions.

## 6. Only Boolean Values in `if` Conditions

Always use boolean values in `if` conditions. Languages like PHP and JS cast non-boolean values to booleans, which requires additional analysis from the programmer and reduces readability.

## 7. Avoiding Negations in Logical Conditions

Negations in conditions require additional analysis. Avoid them where possible by using positive phrasing for variable names or conditions. For example, `if (isUserLoggedOut)` is clearer than `if (!isUserLoggedIn)`.

## 8. Consistency

Adhere to official formatting tools of the stack:
- In **Go**: use standard `gofmt` and standard Go conventions (idiomatic error handling, exported/unexported capitalization).
- In **TypeScript/React**: adhere to `ESLint` and `Prettier` configurations defined in `frontend`.
In files with existing different formatting, maintain the current standard for small changes. Simple auto-formatting changes to a single line are always acceptable.

## 9. Refactor

Adopt the principle: "Leave the campground cleaner than you found it." If you are adding to an existing function and notice code that could be improved, fix it. This enhances readability.

**Important Considerations for Refactoring:**

*   **Common sense**: Refactor existing functionality with common sense.
*   **No refactoring during ASAP tasks.**
*   **Focused Refactoring**: Refactor specific areas, not the entire codebase.
*   **Planned Refactoring**: Every refactor should have a clear goal (simplification, optimization, formatting) and align with the Codebook and Best Practices.
*   **Consultation for Major Refactors**: Consult with the TL and/or PO for any major refactor.
*   **Critical Areas**: Exercise caution when refactoring critical system components.

## 10. Minimizing Risk

*   **Backward Compatibility**: When modifying existing functionalities, ensure changes do not break features currently in use. When modifying common/global code, limit the change's scope to the specific case and ensure old behavior in other contexts.
*   **Modifications in Base Classes**: Prioritize using existing base class mechanisms. Changes to base classes to support derived class processes are rare and should be carefully considered to ensure existing mechanisms cannot handle the task.
*   **Modifications in Existing Core Files**: Avoid modifying classes and files used extensively (e.g., `Form.class`, `baselinker.js`). Such modifications carry higher risk. Always confirm necessity and potential impact on other areas.
*   **Modifications to Existing Functions**: Before changing a function used in multiple places, ensure the change does not break other uses.

## 11. Foolproofing

Code should be "foolproof." This minimizes errors and increases application stability.

*   **Reduced Risk of Errors**: The system should detect and handle errors even if a developer makes a mistake or forgets validation.
*   **Increased Application Stability**: Error-resistant code ensures correct application function, even with unexpected use, making it more reliable and less prone to failures, especially in large or critical systems.

**Examples of Foolproofing:**

*   Enforcing structure of inheriting classes via abstract methods in the base class.
*   Ensuring a class is properly initialized with data before use.
*   If method A depends on method B (e.g., setting a class variable), ensure method B is always called before method A (e.g., by placing method B at the start of method A).
*   Implement error handling at the lowest possible level (e.g., inside the method) to ensure it is never skipped.
*   Filter user input data (e.g., validate types).

## 12. Comments

Comments are largely unnecessary because good code should be fully understandable without them. If comments are used, they must be updated whenever the code changes to avoid misleading information.

Proper variable naming or extracting code into well-named methods often eliminates the need for comments while maintaining readability.

If comments are absolutely necessary, they should be written in English.
