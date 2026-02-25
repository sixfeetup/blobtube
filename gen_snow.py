#!/usr/bin/env python3
"""Generate a 128x128 animated winter snow scene GIF."""

import math
import random
from PIL import Image, ImageDraw

W, H = 128, 128
FRAMES = 60
FPS = 12

random.seed(42)


def lerp_color(c1, c2, t):
    t = max(0, min(1, t))
    return tuple(int(a + (b - a) * t) for a, b in zip(c1, c2))


# Snowflakes: (x, y, speed, size, drift_speed)
snowflakes = []
for _ in range(60):
    snowflakes.append((
        random.randint(0, 127),
        random.randint(-128, 127),
        random.uniform(0.5, 1.8),
        random.randint(1, 3),
        random.uniform(-0.3, 0.3),
    ))

# Stars
stars = [(random.randint(0, 127), random.randint(0, 45), random.uniform(0.1, 0.4)) for _ in range(40)]

# Pine tree positions: (x, height)
trees = [(8, 28), (20, 35), (55, 22), (65, 30), (112, 32), (120, 25)]


def draw_sky(draw, frame):
    for y in range(H):
        t = y / H
        c = lerp_color((15, 15, 50), (40, 40, 80), t)
        draw.line([(0, y), (W, y)], fill=c)


def draw_stars(draw, frame):
    for sx, sy, speed in stars:
        twinkle = (math.sin(frame * speed + sx * 0.5) + 1) / 2
        b = int(twinkle * 200 + 55)
        draw.point((sx, sy), fill=(b, b, min(255, b + 20)))


def draw_moon(draw, frame):
    mx, my = 105, 18
    # Full moon glow
    for r in range(15, 8, -1):
        t = (r - 8) / 7
        c = lerp_color((200, 200, 220), (40, 40, 70), t)
        draw.ellipse([mx - r, my - r, mx + r, my + r], fill=c)
    draw.ellipse([mx - 8, my - 8, mx + 8, my + 8], fill=(220, 220, 235))
    # Crescent shadow
    draw.ellipse([mx - 4, my - 8, mx + 10, my + 8], fill=lerp_color((220, 220, 235), (15, 15, 50), 0.7))


def draw_hills(draw, frame):
    # Background hills with snow
    for x in range(W):
        hill_y = 75 + math.sin(x * 0.03) * 10 + math.sin(x * 0.07 + 2) * 5
        for y in range(int(hill_y), H):
            t = (y - hill_y) / max(1, H - hill_y)
            c = lerp_color((200, 210, 230), (180, 190, 210), t)
            draw.point((x, y), fill=c)

    # Foreground snow ground
    for x in range(W):
        ground_y = 95 + math.sin(x * 0.05) * 3
        for y in range(int(ground_y), H):
            t = (y - ground_y) / max(1, H - ground_y)
            c = lerp_color((220, 225, 240), (200, 205, 220), t)
            draw.point((x, y), fill=c)


def draw_pine_tree(draw, x, height, frame):
    base_y = 95 + int(math.sin(x * 0.05) * 3)
    trunk_h = 5
    # Trunk
    draw.rectangle([x - 1, base_y - trunk_h, x + 1, base_y], fill=(80, 50, 30))

    # Layers of branches with snow
    layers = height // 7
    for i in range(layers):
        t = i / max(layers - 1, 1)
        layer_y = base_y - trunk_h - i * 7
        width = int((1 - t * 0.6) * 10)
        # Green triangle
        draw.polygon([
            (x - width, layer_y),
            (x, layer_y - 8),
            (x + width, layer_y),
        ], fill=(20, 60 + i * 10, 30))
        # Snow on top
        draw.polygon([
            (x - width + 2, layer_y - 1),
            (x, layer_y - 8),
            (x + width - 2, layer_y - 1),
        ], fill=(220, 230, 245))
        # Snow caps
        draw.line([(x - width + 1, layer_y), (x + width - 1, layer_y)], fill=(210, 220, 235))


def draw_cabin(draw, frame):
    cx, cy = 82, 82
    # Cabin body
    draw.rectangle([cx - 12, cy - 10, cx + 12, cy + 10], fill=(100, 60, 30))
    # Darker logs
    for i in range(5):
        y = cy - 8 + i * 4
        draw.line([(cx - 12, y), (cx + 12, y)], fill=(80, 45, 20))
    # Roof
    draw.polygon([(cx - 15, cy - 10), (cx, cy - 22), (cx + 15, cy - 10)], fill=(120, 70, 35))
    # Snow on roof
    draw.polygon([(cx - 14, cy - 11), (cx, cy - 22), (cx + 14, cy - 11)], fill=(215, 225, 240))
    # Door
    draw.rectangle([cx - 3, cy + 2, cx + 3, cy + 10], fill=(70, 40, 15))
    draw.point((cx + 2, cy + 6), fill=(200, 180, 50))
    # Window with warm glow
    glow_pulse = 0.8 + math.sin(frame * 0.15) * 0.2
    glow = tuple(int(v * glow_pulse) for v in (255, 200, 80))
    draw.rectangle([cx + 5, cy - 5, cx + 10, cy], fill=glow)
    draw.rectangle([cx - 10, cy - 5, cx - 5, cy], fill=glow)
    # Window cross
    draw.line([(cx + 5, cy - 2.5), (cx + 10, cy - 2.5)], fill=(80, 50, 25))
    draw.line([(cx + 7.5, cy - 5), (cx + 7.5, cy)], fill=(80, 50, 25))
    draw.line([(cx - 10, cy - 2.5), (cx - 5, cy - 2.5)], fill=(80, 50, 25))
    draw.line([(cx - 7.5, cy - 5), (cx - 7.5, cy)], fill=(80, 50, 25))

    # Chimney
    draw.rectangle([cx + 6, cy - 22, cx + 10, cy - 15], fill=(140, 80, 50))
    # Snow on chimney
    draw.rectangle([cx + 5, cy - 23, cx + 11, cy - 21], fill=(215, 225, 240))


def draw_smoke(draw, frame):
    cx = 90
    start_y = 60
    for i in range(8):
        t = i / 7
        x = cx + math.sin(frame * 0.1 + i * 0.8) * (3 + t * 5)
        y = start_y - i * 4
        r = 2 + t * 3
        alpha = 1 - t
        c = lerp_color((40, 40, 80), (150, 150, 170), alpha * 0.5)
        draw.ellipse([x - r, y - r, x + r, y + r], fill=c)


def draw_snowman(draw, frame):
    sx, sy = 42, 98
    bob = math.sin(frame * 0.1) * 0.5

    # Bottom ball
    draw.ellipse([sx - 7, sy - 5 + bob, sx + 7, sy + 7 + bob], fill=(230, 235, 245))
    # Middle ball
    draw.ellipse([sx - 5, sy - 12 + bob, sx + 5, sy - 2 + bob], fill=(235, 240, 248))
    # Head
    draw.ellipse([sx - 4, sy - 19 + bob, sx + 4, sy - 11 + bob], fill=(240, 242, 250))

    # Eyes
    draw.point((sx - 2, int(sy - 16 + bob)), fill=(20, 20, 20))
    draw.point((sx + 2, int(sy - 16 + bob)), fill=(20, 20, 20))
    # Nose (carrot)
    draw.line([(sx, int(sy - 14 + bob)), (sx + 3, int(sy - 14 + bob))], fill=(240, 140, 40))
    # Mouth dots
    for i in range(3):
        draw.point((sx - 1 + i, int(sy - 12 + bob)), fill=(20, 20, 20))
    # Buttons
    for i in range(3):
        draw.point((sx, int(sy - 9 + i * 3 + bob)), fill=(20, 20, 20))

    # Hat
    hat_y = int(sy - 19 + bob)
    draw.rectangle([sx - 5, hat_y - 1, sx + 5, hat_y], fill=(30, 30, 40))
    draw.rectangle([sx - 3, hat_y - 7, sx + 3, hat_y - 1], fill=(30, 30, 40))

    # Scarf
    draw.line([(sx - 5, int(sy - 11 + bob)), (sx + 5, int(sy - 11 + bob))], fill=(200, 40, 40), width=2)
    draw.line([(sx + 4, int(sy - 11 + bob)), (sx + 6, int(sy - 8 + bob))], fill=(200, 40, 40), width=2)

    # Stick arms
    draw.line([(sx - 5, int(sy - 8 + bob)), (sx - 12, int(sy - 14 + bob))], fill=(80, 50, 30))
    draw.line([(sx + 5, int(sy - 8 + bob)), (sx + 12, int(sy - 14 + bob))], fill=(80, 50, 30))


def draw_snowflakes(draw, frame):
    for start_x, start_y, speed, size, drift in snowflakes:
        y = (start_y + frame * speed * 2) % (H + 30) - 15
        x = (start_x + math.sin(frame * 0.1 + start_x) * 3 + frame * drift) % W
        if 0 <= int(x) < W and 0 <= int(y) < H:
            if size <= 1:
                draw.point((int(x), int(y)), fill=(230, 235, 250))
            else:
                draw.ellipse([x - size / 2, y - size / 2, x + size / 2, y + size / 2],
                             fill=(220, 230, 248))


frames = []
for f in range(FRAMES):
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)

    draw_sky(draw, f)
    draw_stars(draw, f)
    draw_moon(draw, f)
    draw_hills(draw, f)
    draw_smoke(draw, f)
    for tx, th in trees:
        draw_pine_tree(draw, tx, th, f)
    draw_cabin(draw, f)
    draw_snowman(draw, f)
    draw_snowflakes(draw, f)

    frames.append(img)

output_path = "/Users/josh/Code/sixfeetup/blobtube/snow.gif"
frames[0].save(
    output_path,
    save_all=True,
    append_images=frames[1:],
    duration=1000 // FPS,
    loop=0,
)
print(f"Saved {len(frames)} frames to {output_path}")
