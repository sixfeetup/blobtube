#!/usr/bin/env python3
"""Generate a 128x128 animated space scene GIF."""

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


# Stars: (x, y, brightness, twinkle_speed)
stars = [(random.randint(0, 127), random.randint(0, 127), random.uniform(0.3, 1.0), random.uniform(0.1, 0.4)) for _ in range(80)]

# Nebula blobs: (x, y, radius, color)
nebula_blobs = [
    (30, 80, 30, (60, 20, 80)),
    (90, 40, 25, (20, 30, 70)),
    (60, 60, 20, (50, 15, 50)),
    (110, 90, 18, (30, 20, 60)),
]


def draw_background(draw, frame):
    for y in range(H):
        t = y / H
        c = lerp_color((5, 5, 20), (15, 5, 30), t)
        draw.line([(0, y), (W, y)], fill=c)


def draw_nebula(img, frame):
    pixels = img.load()
    for bx, by, br, bc in nebula_blobs:
        pulse = 1.0 + math.sin(frame * 0.05 + bx * 0.1) * 0.15
        r = br * pulse
        for dy in range(int(-r), int(r) + 1):
            for dx in range(int(-r), int(r) + 1):
                dist = math.sqrt(dx * dx + dy * dy)
                if dist < r:
                    px, py = int(bx + dx), int(by + dy)
                    if 0 <= px < W and 0 <= py < H:
                        t = 1 - dist / r
                        alpha = t * t * 0.3
                        existing = pixels[px, py]
                        blended = tuple(min(255, int(existing[i] + bc[i] * alpha)) for i in range(3))
                        pixels[px, py] = blended


def draw_stars(draw, frame):
    for sx, sy, brightness, speed in stars:
        twinkle = (math.sin(frame * speed + sx + sy) + 1) / 2
        b = int(brightness * twinkle * 255)
        if b > 30:
            c = (b, b, min(255, b + 30))
            draw.point((sx, sy), fill=c)
            if b > 180:
                draw.point((sx + 1, sy), fill=(b // 2, b // 2, b // 2))
                draw.point((sx - 1, sy), fill=(b // 2, b // 2, b // 2))


def draw_earth(draw, frame):
    cx, cy = 50, 55
    r = 18
    rotation = frame * 0.03
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            dist = math.sqrt(dx * dx + dy * dy)
            if dist <= r:
                px, py = cx + dx, cy + dy
                # Sphere shading
                norm_x = dx / r
                norm_y = dy / r
                light = max(0, -norm_x * 0.5 + norm_y * -0.3 + 0.6)

                # Simple continent pattern
                angle = math.atan2(dy, dx) + rotation
                lat = norm_y
                pattern = math.sin(angle * 3 + lat * 4) * math.cos(angle * 2 - lat * 3)

                if pattern > 0.1:
                    base = (30, 130, 50)  # Land
                else:
                    base = (20, 60, 180)  # Ocean

                c = tuple(min(255, int(v * light)) for v in base)
                # Atmosphere glow at edges
                edge = dist / r
                if edge > 0.85:
                    atmo_t = (edge - 0.85) / 0.15
                    c = lerp_color(c, (100, 150, 255), atmo_t * 0.6)
                draw.point((px, py), fill=c)


def draw_saturn(draw, frame):
    cx, cy = 100, 30
    r = 10
    # Planet body
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            dist = math.sqrt(dx * dx + dy * dy)
            if dist <= r:
                px, py = cx + dx, cy + dy
                light = max(0.2, (-dx / r * 0.4 + 0.7))
                band = math.sin(dy * 0.8) * 0.15
                base = (210 + int(band * 40), 180 + int(band * 30), 120)
                c = tuple(min(255, int(v * light)) for v in base)
                draw.point((px, py), fill=c)

    # Rings
    ring_tilt = 0.3
    for angle_deg in range(360):
        a = math.radians(angle_deg)
        for ring_r in range(14, 20):
            rx = cx + math.cos(a) * ring_r
            ry = cy + math.sin(a) * ring_r * ring_tilt
            if 0 <= int(rx) < W and 0 <= int(ry) < H:
                # Don't draw ring behind planet
                if math.sin(a) > 0 and abs(rx - cx) < r and abs(ry - cy) < r * ring_tilt:
                    continue
                brightness = 0.5 + math.sin(ring_r * 1.5) * 0.3
                c = tuple(int(v * brightness) for v in (200, 190, 160))
                draw.point((int(rx), int(ry)), fill=c)


def draw_moon(draw, frame):
    mx, my = 18, 20
    r = 8
    for dy in range(-r, r + 1):
        for dx in range(-r, r + 1):
            dist = math.sqrt(dx * dx + dy * dy)
            if dist <= r:
                px, py = mx + dx, my + dy
                light = max(0.2, (dx / r * 0.3 + 0.6))
                # Craters
                crater = 0
                for cx2, cy2, cr in [(2, -2, 2), (-3, 1, 1.5), (1, 3, 1)]:
                    d = math.sqrt((dx - cx2) ** 2 + (dy - cy2) ** 2)
                    if d < cr:
                        crater = 0.15
                base = (180 - int(crater * 100), 180 - int(crater * 100), 170 - int(crater * 100))
                c = tuple(min(255, int(v * light)) for v in base)
                draw.point((px, py), fill=c)


def draw_comet(draw, frame):
    # Comet appears periodically
    cycle = frame % 60
    if cycle > 40:
        return
    t = cycle / 40
    cx = int(W * (1 - t) + 10)
    cy = int(10 + t * 50)

    # Tail
    tail_len = 20
    for i in range(tail_len):
        tt = i / tail_len
        tx = cx + i * 1.2
        ty = cy - i * 0.5
        brightness = int((1 - tt) * 200)
        if brightness > 10 and 0 <= int(tx) < W and 0 <= int(ty) < H:
            draw.point((int(tx), int(ty)), fill=(brightness, brightness, brightness // 2))

    # Head
    draw.ellipse([cx - 2, cy - 2, cx + 2, cy + 2], fill=(255, 255, 200))
    draw.ellipse([cx - 1, cy - 1, cx + 1, cy + 1], fill=(255, 255, 255))


def draw_rocket(draw, frame):
    rx = 75 + math.sin(frame * 0.08) * 5
    ry = 95 + math.cos(frame * 0.06) * 3
    bob = math.sin(frame * 0.15)

    # Body
    draw.rectangle([rx - 2, ry - 6, rx + 2, ry + 4], fill=(200, 200, 210))
    # Nose
    draw.polygon([(rx - 2, ry - 6), (rx, ry - 10), (rx + 2, ry - 6)], fill=(220, 50, 50))
    # Fins
    draw.polygon([(rx - 2, ry + 2), (rx - 5, ry + 5), (rx - 2, ry + 4)], fill=(220, 50, 50))
    draw.polygon([(rx + 2, ry + 2), (rx + 5, ry + 5), (rx + 2, ry + 4)], fill=(220, 50, 50))
    # Window
    draw.ellipse([rx - 1, ry - 4, rx + 1, ry - 2], fill=(100, 180, 255))
    # Flame
    flame_len = 3 + abs(bob) * 3
    draw.polygon([(rx - 2, ry + 4), (rx, ry + 4 + flame_len), (rx + 2, ry + 4)],
                 fill=(255, 200, 50))
    draw.polygon([(rx - 1, ry + 4), (rx, ry + 4 + flame_len * 0.7), (rx + 1, ry + 4)],
                 fill=(255, 255, 150))


def draw_ufo(draw, frame):
    ux = (frame * 1.5 + 20) % (W + 30) - 15
    uy = 108 + math.sin(frame * 0.12) * 4

    # Dome
    draw.ellipse([ux - 4, uy - 6, ux + 4, uy - 1], fill=(150, 200, 150))
    # Saucer
    draw.ellipse([ux - 8, uy - 3, ux + 8, uy + 2], fill=(180, 180, 190))
    # Lights
    for i in range(3):
        lx = ux - 5 + i * 5
        blink = math.sin(frame * 0.4 + i * 2) > 0
        color = (255, 255, 0) if blink else (100, 100, 50)
        draw.ellipse([lx - 1, uy - 1, lx + 1, uy + 1], fill=color)


frames = []
for f in range(FRAMES):
    img = Image.new("RGB", (W, H))
    draw = ImageDraw.Draw(img)

    draw_background(draw, f)
    draw_nebula(img, f)
    draw_stars(draw, f)
    draw_moon(draw, f)
    draw_earth(draw, f)
    draw_saturn(draw, f)
    draw_comet(draw, f)
    draw_rocket(draw, f)
    draw_ufo(draw, f)

    frames.append(img)

output_path = "/Users/josh/Code/sixfeetup/blobtube/space.gif"
frames[0].save(
    output_path,
    save_all=True,
    append_images=frames[1:],
    duration=1000 // FPS,
    loop=0,
)
print(f"Saved {len(frames)} frames to {output_path}")
