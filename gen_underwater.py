#!/usr/bin/env python3
"""Generate a 128x128 animated underwater scene GIF."""

import math
import random
from PIL import Image, ImageDraw

W, H = 128, 128
FRAMES = 60
FPS = 12

random.seed(42)

# Colors
WATER_TOP = (30, 120, 200)
WATER_BOT = (10, 30, 80)
SAND_COLOR = (200, 180, 130)
CORAL_COLORS = [(220, 80, 80), (240, 140, 60), (200, 60, 160), (255, 180, 80)]
KELP_GREEN = (20, 120, 40)
KELP_LIGHT = (40, 160, 60)
BUBBLE_COLOR = (150, 200, 255)

# Fish definitions: (start_x, y, speed, size, color, direction)
fish_list = []
for _ in range(6):
    fish_list.append((
        random.randint(0, 127),
        random.randint(20, 95),
        random.uniform(0.8, 2.0),
        random.randint(4, 8),
        (random.randint(150, 255), random.randint(80, 255), random.randint(50, 255)),
        random.choice([-1, 1]),
    ))

# Bubbles: (x, start_y, speed, size)
bubbles = [(random.randint(10, 118), random.randint(40, 120), random.uniform(0.5, 1.5), random.randint(1, 3)) for _ in range(15)]

# Kelp positions
kelps = [(random.randint(5, 123), random.randint(25, 40)) for _ in range(8)]

# Shells
shells = [(random.randint(10, 118), random.randint(112, 122)) for _ in range(5)]

# Light rays
rays = [(random.randint(10, 118), random.randint(8, 20)) for _ in range(4)]


def lerp_color(c1, c2, t):
    t = max(0, min(1, t))
    return tuple(int(a + (b - a) * t) for a, b in zip(c1, c2))


def draw_water_bg(draw, frame):
    for y in range(H):
        t = y / H
        c = lerp_color(WATER_TOP, WATER_BOT, t)
        draw.line([(0, y), (W, y)], fill=c)


def draw_light_rays(draw, frame):
    for rx, rw in rays:
        sway = math.sin(frame * 0.08 + rx * 0.1) * 5
        x = rx + sway
        for y in range(0, 90):
            t = y / 90
            spread = t * rw
            alpha_t = (1 - t) * 0.15
            c = lerp_color(WATER_TOP, (180, 220, 255), 1 - t)
            bg = lerp_color(WATER_TOP, WATER_BOT, y / H)
            blended = lerp_color(bg, c, alpha_t)
            for dx in range(int(-spread), int(spread) + 1):
                px = int(x + dx)
                if 0 <= px < W:
                    draw.point((px, y), fill=blended)


def draw_sand_floor(draw, frame):
    for y in range(110, H):
        t = (y - 110) / (H - 110)
        c = lerp_color((160, 140, 100), SAND_COLOR, t)
        draw.line([(0, y), (W, y)], fill=c)
    for x in range(0, W, 8):
        sx = x + math.sin(x * 0.3) * 3
        draw.arc([sx, 115, sx + 10, 120], 0, 180, fill=(180, 160, 120))


def draw_coral(draw, frame):
    random.seed(99)
    coral_positions = [(15, 108), (45, 110), (80, 107), (105, 111)]
    for i, (cx, cy) in enumerate(coral_positions):
        color = CORAL_COLORS[i % len(CORAL_COLORS)]
        sway = math.sin(frame * 0.05 + i) * 1
        for branch in range(4):
            angle = -math.pi / 2 + (branch - 1.5) * 0.4 + sway * 0.1
            for seg in range(8):
                t = seg / 7
                bx = cx + math.cos(angle) * seg * 2
                by = cy + math.sin(angle) * seg * 2
                r = 3 - t * 2
                if r > 0:
                    draw.ellipse([bx - r, by - r, bx + r, by + r], fill=color)


def draw_kelp(draw, frame):
    for i, (kx, height) in enumerate(kelps):
        for seg in range(height):
            t = seg / max(height - 1, 1)
            sway = math.sin(frame * 0.1 + i * 0.7 + seg * 0.3) * (3 * t)
            x = kx + sway
            y = 118 - seg * 3
            w = 3 - t * 1.5
            c = lerp_color(KELP_GREEN, KELP_LIGHT, t * 0.5 + math.sin(frame * 0.1 + seg) * 0.2)
            if w > 0:
                draw.ellipse([x - w, y - 1, x + w, y + 1], fill=c)


def draw_fish(draw, frame):
    for start_x, fy, speed, size, color, direction in fish_list:
        x = (start_x + frame * speed * direction) % (W + 20) - 10
        bob = math.sin(frame * 0.2 + start_x) * 2
        y = fy + bob

        if direction > 0:
            draw.ellipse([x - size, y - size // 2, x + size, y + size // 2], fill=color)
            draw.polygon([(x - size, y), (x - size - size // 2, y - size // 2),
                          (x - size - size // 2, y + size // 2)], fill=color)
            draw.ellipse([x + size // 3, y - 2, x + size // 3 + 2, y], fill=(0, 0, 0))
        else:
            draw.ellipse([x - size, y - size // 2, x + size, y + size // 2], fill=color)
            draw.polygon([(x + size, y), (x + size + size // 2, y - size // 2),
                          (x + size + size // 2, y + size // 2)], fill=color)
            draw.ellipse([x - size // 3 - 2, y - 2, x - size // 3, y], fill=(0, 0, 0))

        darker = tuple(max(0, c - 40) for c in color)
        draw.line([(x - size // 2, y), (x + size // 2, y)], fill=darker, width=1)


def draw_jellyfish(draw, frame):
    jx = 95 + math.sin(frame * 0.06) * 10
    jy = 35 + math.sin(frame * 0.08) * 8
    pulse = math.sin(frame * 0.2) * 2

    jelly_color = (200, 150, 255)
    r = 8 + pulse
    draw.ellipse([jx - r, jy - r * 0.6, jx + r, jy + r * 0.4], fill=jelly_color)
    draw.ellipse([jx - r + 2, jy - r * 0.4, jx + r - 2, jy + r * 0.2],
                 fill=(220, 180, 255))

    for t_i in range(5):
        tx = jx - 6 + t_i * 3
        for seg in range(8):
            ty = jy + r * 0.3 + seg * 2.5
            sway = math.sin(frame * 0.15 + t_i + seg * 0.4) * 2
            draw.point((int(tx + sway), int(ty)), fill=(180, 130, 230))


def draw_starfish(draw, frame):
    sx, sy = 30, 117
    color = (230, 140, 50)
    for i in range(5):
        angle = -math.pi / 2 + i * 2 * math.pi / 5
        for seg in range(6):
            t = seg / 5
            px = sx + math.cos(angle) * seg * 1.5
            py = sy + math.sin(angle) * seg * 1.5
            r = 1.5 - t * 0.8
            if r > 0:
                draw.ellipse([px - r, py - r, px + r, py + r], fill=color)
    draw.ellipse([sx - 1, sy - 1, sx + 1, sy + 1], fill=(240, 160, 70))


def draw_bubbles(draw, frame):
    for bx, start_y, speed, size in bubbles:
        y = (start_y - frame * speed) % (H + 20)
        if y < 5 or y > 120:
            continue
        wobble = math.sin(frame * 0.2 + bx) * 1.5
        x = bx + wobble
        draw.ellipse([x - size, y - size, x + size, y + size], outline=BUBBLE_COLOR)
        draw.point((int(x - size * 0.3), int(y - size * 0.3)), fill=(200, 230, 255))


def draw_shells(draw, frame):
    for sx, sy in shells:
        draw.arc([sx - 3, sy - 2, sx + 3, sy + 2], 0, 180, fill=(230, 210, 180))
        draw.arc([sx - 2, sy - 1, sx + 2, sy + 1], 0, 180, fill=(210, 190, 160))


frames = []
for f in range(FRAMES):
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)

    draw_water_bg(draw, f)
    draw_light_rays(draw, f)
    draw_sand_floor(draw, f)
    draw_coral(draw, f)
    draw_kelp(draw, f)
    draw_starfish(draw, f)
    draw_shells(draw, f)
    draw_fish(draw, f)
    draw_jellyfish(draw, f)
    draw_bubbles(draw, f)

    frames.append(img)

output_path = "/Users/josh/Code/sixfeetup/blobtube/underwater.gif"
frames[0].save(
    output_path,
    save_all=True,
    append_images=frames[1:],
    duration=1000 // FPS,
    loop=0,
)
print(f"Saved {len(frames)} frames to {output_path}")
