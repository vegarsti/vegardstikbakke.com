#!/usr/bin/env python3
"""
Generate a default Open Graph image (1200x630px) for social media previews.
"""

from PIL import Image, ImageDraw, ImageFont, ImageOps
import sys
import os

def generate_og_image(output_path, title="Vegard Stikbakke", subtitle="Software Engineer", photo_path="static/me.jpg"):
    # Create image with dimensions 1200x630 (OG image standard)
    width, height = 1200, 630

    # Light mode colors matching your site
    bg_color = (255, 255, 255)  # White background
    text_color = (17, 17, 17)   # Dark text #111
    accent_color = (0, 102, 204)  # Blue link color #0066cc

    # Create image
    img = Image.new('RGB', (width, height), bg_color)
    draw = ImageDraw.Draw(img)

    # Try to use Inter font (similar to site), fall back to SF Pro or Helvetica
    try:
        # Try SF Pro (comes with macOS, similar to Inter)
        title_font = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Display-Semibold.otf", 72)
        subtitle_font = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Display-Regular.otf", 36)
    except:
        try:
            # Fallback to SF Pro Text
            title_font = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Text-Semibold.otf", 72)
            subtitle_font = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Text-Regular.otf", 36)
        except:
            try:
                # Fallback to Helvetica
                title_font = ImageFont.truetype("/System/Library/Fonts/Helvetica.ttc", 72)
                subtitle_font = ImageFont.truetype("/System/Library/Fonts/Helvetica.ttc", 36)
            except:
                # Final fallback
                title_font = ImageFont.load_default()
                subtitle_font = ImageFont.load_default()

    # Add profile photo if it exists
    if os.path.exists(photo_path):
        try:
            # Load and resize photo to square
            photo = Image.open(photo_path)
            photo_size = 180
            photo = ImageOps.fit(photo, (photo_size, photo_size), Image.Resampling.LANCZOS)

            # Position photo on the left side (square, no circular mask)
            photo_x = 80
            photo_y = (height - photo_size) // 2
            img.paste(photo, (photo_x, photo_y))

            # Adjust text position to accommodate photo
            text_start_x = photo_x + photo_size + 60
        except Exception as e:
            print(f"Warning: Could not load photo: {e}")
            text_start_x = 100
    else:
        text_start_x = 100

    # Draw title
    title_y = (height // 2) - 50
    draw.text((text_start_x, title_y), title, fill=text_color, font=title_font)

    # Draw subtitle below title
    subtitle_y = title_y + 90
    draw.text((text_start_x, subtitle_y), subtitle, fill=text_color, font=subtitle_font)

    # Save image
    img.save(output_path, 'PNG', optimize=True)
    print(f"✓ Generated OG image: {output_path}")

if __name__ == "__main__":
    output = sys.argv[1] if len(sys.argv) > 1 else "static/og-image.png"
    generate_og_image(output)
