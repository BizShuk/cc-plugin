#!/bin/sh
"exec" "python3" "$0" "$@"
# -*- coding: utf-8 -*-
import json
import os
import subprocess
import sys

def render_quota(q, total_window_sec, width=10, fill_char="█"):
    if not q:
        return "N/A"
    rem = q.get("remaining_fraction", 1.0)
    used_pct = max(0.0, min(100.0, (1.0 - rem) * 100.0))
    sec = q.get("reset_in_seconds", total_window_sec)
    elapsed_pct = max(0.0, min(100.0, (total_window_sec - sec) / total_window_sec * 100.0))
    
    filled = max(0, min(width, round(used_pct / 100.0 * width)))
    cursor = max(0, min(width - 1, int(elapsed_pct / 100.0 * width)))
    
    chars = []
    for i in range(width):
        if i == cursor:
            chars.append("│")
        elif i < filled:
            chars.append(fill_char)
        else:
            chars.append("░")
    bar = "".join(chars)
    return f"[{bar}] {used_pct:.1f}%"

def format_tokens(n):
    if n is None:
        return "0"
    if n >= 1_000_000:
        return f"{n/1_000_000:.1f}M"
    if n >= 1_000:
        return f"{n/1_000:.1f}k"
    return str(n)

def get_git_info(cwd):
    try:
        branch = subprocess.check_output(
            ["git", "-C", cwd, "rev-parse", "--abbrev-ref", "HEAD"],
            stderr=subprocess.DEVNULL,
            timeout=0.3
        ).decode().strip()
        if not branch:
            return ""
        status = subprocess.check_output(
            ["git", "-C", cwd, "status", "--porcelain"],
            stderr=subprocess.DEVNULL,
            timeout=0.3
        ).decode().strip()
        dirty = "*" if status else ""
        return f"⎇ {branch}{dirty}"
    except Exception:
        return ""

def main():
    try:
        raw_input = sys.stdin.read()
        if not raw_input.strip():
            return
        data = json.loads(raw_input)
    except Exception:
        return

    term_width = data.get("terminal_width") or 80

    # 1. Model Info
    model_info = data.get("model", {})
    model_name = model_info.get("display_name") or model_info.get("id") or "Gemini"
    effort = model_info.get("effort")
    effort_str = f" [{effort}]" if effort else ""

    is_claude = "claude" in model_name.lower() or "sonnet" in model_name.lower() or "opus" in model_name.lower()

    # 2. Workspace & Git
    cwd = data.get("cwd") or data.get("workspace", {}).get("current_dir") or os.getcwd()
    home = os.path.expanduser("~")
    short_cwd = cwd.replace(home, "~") if cwd.startswith(home) else cwd
    git_info = get_git_info(cwd)

    # 3. Context & Tokens
    ctx = data.get("context_window", {})
    ctx_pct = ctx.get("used_percentage", 0)
    total_in = ctx.get("total_input_tokens", 0)
    total_out = ctx.get("total_output_tokens", 0)
    cur_usage = ctx.get("current_usage") or {}
    cached_tokens = cur_usage.get("cache_read_input_tokens", 0)

    # 4. Quotas: Gemini and Claude (3p)
    quotas = data.get("quota", {})
    
    # Gemini (7d & 5h)
    q_gem_w = quotas.get("gemini-weekly")
    q_gem_5 = quotas.get("gemini-5h")
    gem_w_str = render_quota(q_gem_w, 7 * 86400, width=10, fill_char="█")
    gem_5_str = render_quota(q_gem_5, 5 * 3600, width=10, fill_char="▓")

    # Claude / 3p (7d & 5h)
    q_cla_w = quotas.get("3p-weekly")
    q_cla_5 = quotas.get("3p-5h")
    cla_w_str = render_quota(q_cla_w, 7 * 86400, width=10, fill_char="█")
    cla_5_str = render_quota(q_cla_5, 5 * 3600, width=10, fill_char="▓")

    # ANSI Colors
    C_RESET = "\033[0m"
    C_DIM = "\033[2m"
    C_BLUE = "\033[1;34m"
    C_CYAN = "\033[1;36m"
    C_YELLOW = "\033[1;33m"
    C_MAGENTA = "\033[1;35m"
    C_WHITE = "\033[1;37m"
    SEP = f"{C_DIM} | {C_RESET}"

    # Line 1: Model | CWD | Tokens | Context
    tokens_str = f"{format_tokens(total_in + total_out)} ({format_tokens(total_in)}, {format_tokens(cached_tokens)}, {format_tokens(total_out)})"
    line1_parts = [
        f"{C_BLUE}{model_name}{C_RESET}{C_DIM}{effort_str}{C_RESET}",
        f"{C_WHITE}{short_cwd}{C_RESET}",
        f"{C_DIM}{tokens_str}{C_RESET}",
        f"{C_WHITE}Context: {ctx_pct:.1f}%{C_RESET}"
    ]
    line1 = SEP.join(line1_parts)

    # Line 2: Weekly & 5h limits for Gemini and Claude
    gem_section = f"{C_CYAN}Gemini W: {gem_w_str}{C_RESET} {C_YELLOW}5h: {gem_5_str}{C_RESET}"
    cla_section = f"{C_CYAN}Claude W: {cla_w_str}{C_RESET} {C_YELLOW}5h: {cla_5_str}{C_RESET}"

    # Decide layout based on terminal width and active model
    line2_parts = []
    if term_width >= 115 and q_gem_w and q_cla_w:
        # Wide screen: Show both Gemini and Claude, active model first
        if is_claude:
            line2_parts.extend([cla_section, gem_section])
        else:
            line2_parts.extend([gem_section, cla_section])
        if git_info:
            line2_parts.append(f"{C_MAGENTA}{git_info}{C_RESET}")
    else:
        # Compact / Narrow screen: Show active model
        if is_claude:
            active_section = f"{C_CYAN}Claude W: {cla_w_str}{C_RESET} {C_YELLOW}5h: {cla_5_str}{C_RESET}"
        else:
            active_section = f"{C_CYAN}Gemini W: {gem_w_str}{C_RESET} {C_YELLOW}5h: {gem_5_str}{C_RESET}"
        line2_parts.append(active_section)
        if git_info:
            line2_parts.append(f"{C_MAGENTA}{git_info}{C_RESET}")

    line2 = SEP.join(line2_parts)

    print(f"{line1}\n{line2}")

if __name__ == "__main__":
    main()
