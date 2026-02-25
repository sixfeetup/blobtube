#!/usr/bin/env python3
"""Generate a 128x128 animated campfire night scene GIF."""

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


# Stars
stars = [(random.randint(0, 127), random.randint(0, 50), random.uniform(0.1, 0.4)) for _ in range(50)]

# Fireflies
fireflies = [(random.randint(0, 127), random.randint(30, 90), random.uniform(0.1, 0.3), random.uniform(0.5, 1.5)) for _ in range(12)]

# Sparks from fire
sparks = [(random.uniform(-5, 5), random.uniform(2, 8), random.uniform(0.5, 1.5)) for _ in range(15)]


def draw_sky(draw, frame):
    for y in range(H):
        t = y / H
        c = lerp_color((5, 5, 25), (15, 15, 40), t)
        draw.line([(0, y), (W, y)], fill=c)


def draw_stars(draw, frame):
    for sx, sy, speed in stars:
        twinkle = (math.sin(frame * speed + sx + sy) + 1) / 2
        b = int(twinkle * 180 + 40)
        draw.point((sx, sy), fill=(b, b, min(255, b + 15)))


def draw_crescent_moon(draw, frame):
    mx, my = 20, 15
    draw.ellipse([mx - 7, my - 7, mx + 7, my + 7], fill=(230, 230, 210))
    draw.ellipse([mx - 3, my - 8, mx + 9, my + 6], fill=(5, 5, 25))


def draw_mountains(draw, frame):
    # Far mountains
    for x in range(W):
        y1 = 55 + math.sin(x * 0.02) * 15 + math.sin(x * 0.05 + 1) * 8
        for y in range(int(y1), H):
            draw.point((x, y), fill=(15, 20, 35))

    # Near mountains
    for x in range(W):
        y1 = 70 + math.sin(x * 0.03 + 2) * 10 + math.sin(x * 0.06) * 5
        for y in range(int(y1), H):
            draw.point((x, y), fill=(10, 15, 25))


def draw_ground(draw, frame):
    for x in range(W):
        ground_y = 90 + math.sin(x * 0.04) * 2
        for y in range(int(ground_y), H):
            t = (y - ground_y) / max(1, H - ground_y)
            c = lerp_color((25, 35, 20), (20, 28, 15), t)
            draw.point((x, y), fill=c)


def draw_tree_silhouettes(draw, frame):
    tree_positions = [(5, 60), (15, 55), (25, 58), (100, 57), (110, 53), (120, 60)]
    for tx, ty in tree_positions:
        # Trunk
        draw.rectangle([tx - 1, ty, tx + 1, ty + 30], fill=(8, 10, 15))
        # Canopy triangle
        for i in range(4):
            layer_y = ty - i * 5
            width = 8 - i * 1.5
            draw.polygon([
                (tx - width, layer_y),
                (tx, layer_y - 7),
                (tx + width, layer_y),
            ], fill=(8, 12, 18))


def draw_tent(draw, frame):
    tx, ty = 100, 85
    draw.polygon([(tx - 12, ty + 8), (tx, ty - 8), (tx + 12, ty + 8)], fill=(60, 50, 40))
    draw.polygon([(tx - 10, ty + 8), (tx, ty - 6), (tx + 10, ty + 8)], fill=(70, 60, 50))
    # Door flap
    draw.polygon([(tx - 2, ty + 8), (tx, ty), (tx + 2, ty + 8)], fill=(50, 40, 30))


def draw_fire_glow(img, frame):
    """Add warm glow around the fire."""
    pixels = img.load()
    fire_x, fire_y = 64, 88
    glow_pulse = 0.8 + math.sin(frame * 0.3) * 0.2
    glow_r = 35 * glow_pulse

    for dy in range(int(-glow_r), int(glow_r) + 1):
        for dx in range(int(-glow_r), int(glow_r) + 1):
            dist = math.sqrt(dx * dx + dy * dy)
            if dist < glow_r:
                px, py = fire_x + dx, fire_y + dy
                if 0 <= px < W and 0 <= py < H:
                    t = 1 - dist / glow_r
                    alpha = t * t * 0.25
                    glow_color = (255, 140, 30)
                    existing = pixels[px, py]
                    blended = tuple(min(255, int(existing[i] + glow_color[i] * alpha)) for i in range(3))
                    pixels[px, py] = blended


def draw_logs(draw, frame):
    cx, cy = 64, 95
    # Log 1
    draw.ellipse([cx - 10, cy - 2, cx - 2, cy + 3], fill=(60, 35, 15))
    draw.ellipse([cx + 2, cy - 2, cx + 10, cy + 3], fill=(55, 30, 12))
    # Cross log
    draw.line([(cx - 8, cy + 2), (cx + 8, cy - 2)], fill=(65, 38, 18), width=3)
    # Stones around fire
    stone_color = (80, 80, 85)
    for angle_deg in range(0, 360, 30):
        a = math.radians(angle_deg)
        sx = cx + math.cos(a) * 12
        sy = cy + math.sin(a) * 5
        draw.ellipse([sx - 2, sy - 1.5, sx + 2, sy + 1.5], fill=stone_color)


def draw_fire(draw, frame):
    cx, cy = 64, 88

    # Flames - multiple overlapping shapes
    for i in range(5):
        phase = frame * 0.3 + i * 1.3
        sway = math.sin(phase) * 3
        height = 12 + math.sin(phase * 1.5) * 4
        width = 5 + math.sin(phase * 0.7) * 2

        fx = cx + sway + (i - 2) * 2
        fy = cy

        # Outer flame (red/orange)
        draw.polygon([
            (fx - width, fy),
            (fx + math.sin(phase + 1) * 2, fy - height),
            (fx + width, fy),
        ], fill=(220, 80 + i * 10, 20))

        # Inner flame (yellow)
        inner_h = height * 0.6
        inner_w = width * 0.5
        draw.polygon([
            (fx - inner_w, fy),
            (fx + math.sin(phase + 2) * 1, fy - inner_h),
            (fx + inner_w, fy),
        ], fill=(255, 200, 50))

    # Bright core
    draw.ellipse([cx - 3, cy - 3, cx + 3, cy + 1], fill=(255, 240, 150))


def draw_sparks(draw, frame):
    cx, cy = 64, 85
    for dx, speed, phase in sparks:
        t = ((frame * 0.08 + phase) % 1.0)
        if t > 0.7:
            continue
        sx = cx + dx + math.sin(frame * 0.2 + phase * 5) * 3
        sy = cy - t * 30
        brightness = (1 - t / 0.7)
        c = (int(255 * brightness), int(150 * brightness), int(30 * brightness))
        draw.point((int(sx), int(sy)), fill=c)


def draw_log_seats(draw, frame):
    # Two log seats
    draw.ellipse([40, 98, 52, 103], fill=(70, 45, 20))
    draw.ellipse([76, 99, 88, 104], fill=(65, 40, 18))


def draw_fireflies(draw, frame):
    for fx, fy, speed, phase in fireflies:
        x = fx + math.sin(frame * speed + phase * 10) * 10
        y = fy + math.cos(frame * speed * 0.7 + phase * 5) * 8
        brightness = (math.sin(frame * speed * 2 + phase * 3) + 1) / 2
        if brightness > 0.6 and 0 <= int(x) < W and 0 <= int(y) < H:
            b = int(brightness * 200)
            draw.point((int(x), int(y)), fill=(b, b, b // 3))


frames = []
for f in range(FRAMES):
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)

    draw_sky(draw, f)
    draw_stars(draw, f)
    draw_crescent_moon(draw, f)
    draw_mountains(draw, f)
    draw_ground(draw, f)
    draw_tree_silhouettes(draw, f)
    draw_tent(draw, f)
    draw_fire_glow(img, f)
    draw_logs(draw, f)
    draw_log_seats(draw, f)
    draw_fire(draw, f)
    draw_sparks(draw, f)
    draw_fireflies(draw, f)

    frames.append(img)

output_path = "/Users/josh/Code/sixfeetup/blobtube/campfire.gif"
frames[0].save(
    output_path,
    save_all=True,
    append_images=frames[1:],
    duration=1000 // FPS,
    loop=0,
)
print(f"Saved {len(frames)} frames to {output_path}")
