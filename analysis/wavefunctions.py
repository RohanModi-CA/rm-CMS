#!/usr/bin/env python3
"""
wavefunctions.py - Quantum mechanical wavefunction analysis.

Computes wave functions and probability densities for simple 1-D quantum
systems (particle-in-a-box and harmonic oscillator) and renders a basic
ASCII plot of the probability density.

Usage:
    python wavefunctions.py [n] [L]

    n  - quantum number  (default: 1)
    L  - box length / a.u. (default: 1.0)
"""

import math
import os
import sys


# ---------------------------------------------------------------------------
# Physics helpers
# ---------------------------------------------------------------------------

def particle_in_box_wavefunction(n, L, x):
    """Return the nth normalised wavefunction for a particle in a 1-D box.

    Parameters
    ----------
    n : int   Quantum number (>= 1).
    L : float Box length (must be > 0).
    x : float Position.

    Raises
    ------
    ValueError if n < 1 or L <= 0.
    """
    if L <= 0:
        raise ValueError("Box length L must be positive (got {}).".format(L))
    if n < 1:
        raise ValueError("Quantum number n must be >= 1 (got {}).".format(n))
    return math.sqrt(2.0 / L) * math.sin(n * math.pi * x / L)


def probability_density(n, L, x):
    """Return |psi_n(x)|^2 for a particle in a 1-D box."""
    psi = particle_in_box_wavefunction(n, L, x)
    return psi * psi


def normalize(values, dx):
    """Normalise *values* so that sum(values)*dx == 1.

    Parameters
    ----------
    values : list[float]
    dx     : float  Grid spacing.  Must be > 0.

    Raises
    ------
    ValueError if dx <= 0 or the total probability is zero.
    """
    if dx <= 0:
        raise ValueError(
            "Grid spacing dx must be positive (got {}).".format(dx)
        )
    total = sum(values) * dx
    if total == 0.0:
        raise ValueError(
            "Cannot normalise: integrated probability is zero."
        )
    return [v / total for v in values]


# ---------------------------------------------------------------------------
# Terminal helpers
# ---------------------------------------------------------------------------

def get_terminal_columns(fallback=80):
    """Return the terminal width, falling back to *fallback* when:
    - the environment variable COLUMNS is absent, empty, or '0'  (common
      on Windows Git Bash, which does not export COLUMNS reliably), or
    - os.get_terminal_size() raises an OSError (non-interactive context).

    This prevents the 'division by zero' error that arises on Windows Git
    Bash when code blindly uses ``int(os.environ['COLUMNS'])`` as a divisor.
    """
    # 1. Try the environment variable first (most portable).
    columns_env = os.environ.get("COLUMNS", "")
    if columns_env:
        try:
            columns = int(columns_env)
            if columns > 0:
                return columns
        except ValueError:
            pass  # not a valid integer – fall through

    # 2. Ask the OS for the real terminal size.
    try:
        size = os.get_terminal_size()
        if size.columns > 0:
            return size.columns
    except OSError:
        pass  # not attached to a terminal (e.g. piped output)

    return fallback


# ---------------------------------------------------------------------------
# ASCII plot
# ---------------------------------------------------------------------------

def ascii_plot(xs, ys, width=None, height=20, label="P(x)"):
    """Print a simple ASCII bar chart of *ys* vs *xs*.

    Parameters
    ----------
    xs     : list[float]
    ys     : list[float]
    width  : int | None  Plot width in characters.  Defaults to terminal width.
    height : int         Number of rows in the plot.
    label  : str         Y-axis label.
    """
    if width is None:
        width = get_terminal_columns()

    # Guard: ensure width is sensible.
    if width <= 0:
        width = 80

    y_max = max(ys) if ys else 1.0
    if y_max == 0.0:
        y_max = 1.0  # avoid division by zero when all values are zero

    # Resample ys to *width* buckets.
    n = len(ys)
    if n == 0:
        print("(no data)")
        return

    bucket_size = n / width  # float; intentionally not integer division
    if bucket_size <= 0:
        bucket_size = 1.0

    buckets = []
    for col in range(width):
        start = int(col * bucket_size)
        end = int((col + 1) * bucket_size)
        end = max(end, start + 1)
        end = min(end, n)
        buckets.append(sum(ys[start:end]) / (end - start))

    print("{} {}".format(label, " " * (width - len(label) - 1)))
    for row in range(height, 0, -1):
        threshold = (row / height) * y_max
        line = ""
        for val in buckets:
            line += "#" if val >= threshold else " "
        print("|{}|".format(line))
    print("+" + "-" * width + "+")
    left_label = "{:.3f}".format(xs[0])
    right_label = "x{:.3f}".format(xs[-1])
    inner_space = width - len(left_label) - len(right_label)
    print(" {}{}{}".format(left_label, " " * max(inner_space, 1), right_label))


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------

def main():
    # Parse optional command-line arguments.
    try:
        n = int(sys.argv[1]) if len(sys.argv) > 1 else 1
    except ValueError:
        print("Error: quantum number n must be an integer.", file=sys.stderr)
        sys.exit(1)

    try:
        L = float(sys.argv[2]) if len(sys.argv) > 2 else 1.0
    except ValueError:
        print("Error: box length L must be a number.", file=sys.stderr)
        sys.exit(1)

    if n < 1:
        print("Error: quantum number n must be >= 1.", file=sys.stderr)
        sys.exit(1)

    if L <= 0:
        print("Error: box length L must be positive.", file=sys.stderr)
        sys.exit(1)

    # Build the probability-density grid.
    num_points = get_terminal_columns() * 4  # 4 samples per column
    if num_points <= 0:
        num_points = 320

    dx = L / num_points

    xs = [i * dx for i in range(num_points + 1)]
    probs = [probability_density(n, L, x) for x in xs]

    try:
        normalised = normalize(probs, dx)
    except ValueError as exc:
        print("Error during normalisation: {}".format(exc), file=sys.stderr)
        sys.exit(1)

    print("Particle-in-a-box  n={}, L={:.3f}".format(n, L))
    ascii_plot(xs, normalised, label="|psi_{}|^2".format(n))


if __name__ == "__main__":
    main()
