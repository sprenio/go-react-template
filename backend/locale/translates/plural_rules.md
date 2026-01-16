# Plural Rules Reference (CLDR)

This document summarizes **plural categories and rules** used by **Unicode CLDR**, which is also what libraries like **go-i18n**, **ICU**, **gettext**, and many frontend frameworks rely on.

Plural categories:

* `zero`
* `one`
* `two`
* `few`
* `many`
* `other`

> ⚠️ Important
>
> * `other` is **NOT a fallback**. If a language requires `few` or `many`, they **must exist**.
> * Most languages use only a subset of these categories.

---

## Category meanings (generic)

| Category | Typical meaning                   |
| -------- | --------------------------------- |
| `zero`   | exactly 0                         |
| `one`    | exactly 1                         |
| `two`    | exactly 2                         |
| `few`    | small numbers (language-specific) |
| `many`   | large numbers (language-specific) |
| `other`  | fallback / decimals               |

---

## Languages with **NO plural distinction**

### Categories required

```
other
```

### Languages

* Japanese (`ja`)
* Chinese (`zh`)
* Korean (`ko`)
* Thai (`th`)
* Vietnamese (`vi`)
* Indonesian (`id`)
* Malay (`ms`)
* Persian (`fa`)
* Turkish (`tr`)

---

## Languages with **2 plural forms** (`one`, `other`)

### Rule

```
one   → n = 1
other → everything else
```

### Languages

* English (`en`)
* German (`de`)
* Dutch (`nl`)
* Swedish (`sv`)
* Danish (`da`)
* Norwegian (`no`)
* Finnish (`fi`)
* Estonian (`et`)
* Spanish (`es`)
* Portuguese (`pt`)
* Italian (`it`)
* Greek (`el`)
* Hebrew (`he`)

⚠️ French note:

* French (`fr`) uses `one` for **0 and 1**

---

## Languages with **3 plural forms** (`one`, `few`, `other`)

### Rule (example: Slovak, Czech)

```
one → n = 1
few → n = 2–4
other → everything else
```

### Languages

* Slovak (`sk`)
* Czech (`cs`)

---

## Languages with **3 plural forms** (`one`, `two`, `other`)

### Rule (example: Arabic subset)

```
one → n = 1
two → n = 2
other → everything else
```

### Languages

* Filipino (`fil`)
* Tagalog (`tl`)

---

## Slavic languages (`one`, `few`, `many`, `other`)

### Rule (generic Slavic)

```
one  → n mod 10 = 1 and n mod 100 ≠ 11
few  → n mod 10 = 2–4 and n mod 100 ≠ 12–14
many → n mod 10 = 0 or 5–9 or n mod 100 = 11–14
other → decimals
```

### Languages

* Polish (`pl`)
* Russian (`ru`)
* Ukrainian (`uk`)
* Belarusian (`be`)
* Serbian (`sr`)
* Croatian (`hr`)
* Bosnian (`bs`)

⚠️ These languages **require `many`**.

---

## Arabic (`zero`, `one`, `two`, `few`, `many`, `other`)

### Rule

```
zero → n = 0
one → n = 1
two → n = 2
few → n mod 100 = 3–10
many → n mod 100 = 11–99
other → decimals
```

### Language

* Arabic (`ar`)

⚠️ Arabic has the **most complex plural rules**.

---

## Irish (`one`, `two`, `few`, `many`, `other`)

### Rule

```
one → n = 1
two → n = 2
few → n = 3–6
many → n = 7–10
other → everything else
```

### Language

* Irish (`ga`)

---

## Welsh (`zero`, `one`, `two`, `few`, `many`, `other`)

### Rule

```
zero → n = 0
one → n = 1
two → n = 2
few → n = 3
many → n = 6
other → everything else
```

### Language

* Welsh (`cy`)

---

## Best Practices

### ✔ Always follow CLDR rules

Do not invent your own plural logic.

### ✔ Slavic languages

Always provide:

```
one
few
many
```

### ✔ Safe universal template

```json
{
  "one": "...",
  "few": "...",
  "many": "...",
  "other": "..."
}
```

Unused forms will be ignored by languages that don’t need them.

---

## References

* Unicode CLDR Plural Rules
* [https://cldr.unicode.org](https://cldr.unicode.org)
* go-i18n documentation

---

*End of document*
