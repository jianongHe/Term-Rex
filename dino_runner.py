import curses
import time
import random
import math

# --- Configuration Constants ---
FPS = 24
TARGET_FRAME_TIME = 1.0 / FPS

# GROUND_CHAR = "_" # Replaced by a more visible one below
# DINO_CHARS = ["#", "@"]  # Removed, using DINO_SPRITES
# OBSTACLE_CHARS = ["|", "#"] # Removed, using OBSTACLE_SPRITES

GROUND_CHAR = "─" # Using a line character for better ground visuals

# --- ASCII Art Definitions ---
DINO_SPRITES = [
    [ # Frame 1
        "  @@ ",
        " @#@ ",
        "@##@ ",
        " ### ",
        " @ @ "
    ],
    [ # Frame 2
        "  @@ ",
        " @#@ ",
        "@##@ ",
        " ### ",
        "@ @  "
    ],
    [ # Frame 3
        "  @@ ",
        " @#@ ",
        "@##@ ",
        " ### ",
        " @ @ "
    ],
    [ # Frame 4
        "  @@ ",
        " @#@ ",
        "@##@ ",
        " ### ",
        "  @ @"
    ]
]

OBSTACLE_SPRITES = [
    [ # Small Cactus
        " # ",
        "###",
        " # "
    ],
    [ # Medium Cactus
        "  #  ",
        " ### ",
        "  #  ",
        "  #  ",
        "  #  "
    ],
    [ # Double Cactus
        " #  # ",
        "### ###",
        " #  # "
    ],
    [ # Tall Cactus
        "  #  ",
        "  #  ",
        " ### ",
        "  #  "
    ]
]

# Colors (Pair Numbers)
COLOR_PAIR_DAY_DINO = 1
COLOR_PAIR_DAY_OBSTACLE = 2
COLOR_PAIR_DAY_GROUND_TEXT = 3
COLOR_PAIR_NIGHT_DINO = 4
COLOR_PAIR_NIGHT_OBSTACLE = 5
COLOR_PAIR_NIGHT_GROUND_TEXT = 6

# Game Dimensions
GROUND_Y_OFFSET = 1  # Visual ground line is 1 unit above the screen bottom.
# Sprites' bottoms align with `ground_line_y`.
DINO_X_POS = 5

# Jump Mechanics
JUMP_HEIGHT = 5
JUMP_DURATION_FRAMES = int(FPS * 0.70) # Slightly faster jump

# Difficulty
INITIAL_OBSTACLE_SPEED = 1.2
INITIAL_OBSTACLE_SPAWN_INTERVAL_SEC = 2.3 # Adjusted
SPEED_INCREASE_PER_STAGE = 0.25
SPAWN_DECREASE_PER_STAGE_SEC = 0.2
STAGE_DURATION_SEC = 30
DAY_NIGHT_CYCLE_SEC = 60
MIN_OBSTACLE_SPAWN_INTERVAL_SEC = 0.7

HIGH_SCORE_FILE = "dino_highscore.txt"

class Player:
    def __init__(self, x, ground_y_coord): # ground_y_coord is where the feet land
        self.x = x # x-coordinate of the left of the sprite
        self.sprites = DINO_SPRITES
        self.sprite_idx = 0

        current_sprite_lines = self.sprites[self.sprite_idx]
        self.height = len(current_sprite_lines)
        self.width = max(len(line) for line in current_sprite_lines) if self.height > 0 else 0

        # y_top is the y-coordinate for the TOP-MOST line of the sprite
        # such that its bottom aligns with ground_y_coord
        self.base_y_top = ground_y_coord - self.height + 1
        self.y_top = self.base_y_top

        self.is_jumping = False
        self.jump_frame_count = 0
        self.animation_timer = 0
        self.animation_speed = FPS // 6 # Animation speed for sprites

    def jump(self):
        if not self.is_jumping:
            self.is_jumping = True
            self.jump_frame_count = 0

    def update(self):
        self.animation_timer +=1
        if self.animation_timer >= self.animation_speed:
            self.sprite_idx = (self.sprite_idx + 1) % len(self.sprites)
            # Assuming all dino animation sprites have the same dimensions for simplicity.
            # If they could change, self.height, self.width, and self.base_y_top
            # would need to be updated here.
            self.animation_timer = 0

        if self.is_jumping:
            self.jump_frame_count += 1
            progress = self.jump_frame_count / JUMP_DURATION_FRAMES
            # Parabolic jump: y_offset_from_base is displacement upwards from base_y_top
            y_offset_from_base = -4 * JUMP_HEIGHT * (progress - 0.5)**2 + JUMP_HEIGHT
            self.y_top = self.base_y_top - int(round(y_offset_from_base))

            if self.jump_frame_count >= JUMP_DURATION_FRAMES:
                self.is_jumping = False
                self.y_top = self.base_y_top # Ensure back on ground correctly

    def get_current_sprite_lines(self):
        return self.sprites[self.sprite_idx]

    def get_rect(self): # Collision rectangle: (left_x, top_y, width, height)
        return (self.x, self.y_top, self.width, self.height)

class Obstacle:
    def __init__(self, x_coord, ground_y_coord, sprite_lines, speed):
        self.sprite_lines = sprite_lines
        self.height = len(self.sprite_lines)
        self.width = max(len(line) for line in self.sprite_lines) if self.height > 0 else 0

        self.x = x_coord # x for the leftmost part of the sprite
        # y_top is the y-coordinate for the TOP-MOST line of the sprite
        # such that its bottom aligns with ground_y_coord
        self.y_top = ground_y_coord - self.height + 1
        self.speed = speed

    def update(self):
        self.x -= self.speed

    def is_offscreen(self):
        return self.x + self.width < 0 # Check if entire sprite is offscreen

    def get_rect(self): # Collision rectangle: (left_x, top_y, width, height)
        return (int(round(self.x)), self.y_top, self.width, self.height)

class Game:
    def __init__(self, stdscr):
        self.stdscr = stdscr
        self.height, self.width = stdscr.getmaxyx()

        self.game_state = "START"
        self.score = 0
        self.high_score = self._load_high_score()
        self.frame_count = 0

        # This is the Y coordinate where the bottom of sprites should align (e.g., feet of dino)
        self.ground_line_y = self.height - 1 - GROUND_Y_OFFSET

        self.player = Player(DINO_X_POS, self.ground_line_y)
        self.obstacles = []

        self.current_obstacle_speed = INITIAL_OBSTACLE_SPEED
        self.current_obstacle_spawn_interval_frames = int(INITIAL_OBSTACLE_SPAWN_INTERVAL_SEC * FPS)
        self.last_obstacle_spawn_frame = 0

        self.stage_timer_frames = 0
        self.day_night_timer_frames = 0
        self.is_day_mode = True

        self._init_colors()

    def _load_high_score(self):
        try:
            with open(HIGH_SCORE_FILE, "r") as f:
                return int(f.read().strip())
        except (FileNotFoundError, ValueError):
            return 0

    def _save_high_score(self):
        try:
            with open(HIGH_SCORE_FILE, "w") as f:
                f.write(str(self.high_score))
        except IOError:
            pass

    def _init_colors(self):
        curses.start_color()
        curses.use_default_colors()
        curses.init_pair(COLOR_PAIR_DAY_DINO, curses.COLOR_GREEN, -1)
        curses.init_pair(COLOR_PAIR_DAY_OBSTACLE, curses.COLOR_RED, -1)
        curses.init_pair(COLOR_PAIR_DAY_GROUND_TEXT, curses.COLOR_WHITE, -1) # Or your terminal's default fg/bg
        curses.init_pair(COLOR_PAIR_NIGHT_DINO, curses.COLOR_CYAN, curses.COLOR_BLUE)
        curses.init_pair(COLOR_PAIR_NIGHT_OBSTACLE, curses.COLOR_MAGENTA, curses.COLOR_BLUE)
        curses.init_pair(COLOR_PAIR_NIGHT_GROUND_TEXT, curses.COLOR_YELLOW, curses.COLOR_BLUE)

    def _get_color_pair(self, element_type):
        if self.is_day_mode:
            if element_type == "dino": return curses.color_pair(COLOR_PAIR_DAY_DINO)
            if element_type == "obstacle": return curses.color_pair(COLOR_PAIR_DAY_OBSTACLE)
            return curses.color_pair(COLOR_PAIR_DAY_GROUND_TEXT)
        else: # Night mode
            if element_type == "dino": return curses.color_pair(COLOR_PAIR_NIGHT_DINO)
            if element_type == "obstacle": return curses.color_pair(COLOR_PAIR_NIGHT_OBSTACLE)
            return curses.color_pair(COLOR_PAIR_NIGHT_GROUND_TEXT)

    def _reset_game(self):
        self.score = 0
        self.frame_count = 0
        self.player = Player(DINO_X_POS, self.ground_line_y) # Use self.ground_line_y
        self.obstacles = []
        self.current_obstacle_speed = INITIAL_OBSTACLE_SPEED
        self.current_obstacle_spawn_interval_frames = int(INITIAL_OBSTACLE_SPAWN_INTERVAL_SEC * FPS)
        self.last_obstacle_spawn_frame = 0
        self.stage_timer_frames = 0
        self.day_night_timer_frames = 0
        self.is_day_mode = True
        self.game_state = "PLAYING"
        # Explicitly set background for day mode on reset too
        bg_color_to_set = curses.color_pair(COLOR_PAIR_DAY_GROUND_TEXT) if self.is_day_mode else curses.color_pair(COLOR_PAIR_NIGHT_GROUND_TEXT)
        try:
            self.stdscr.bkgd(' ', bg_color_to_set)
        except curses.error: # Some terminals might not like changing bkgd for default colors
            if self.is_day_mode: self.stdscr.bkgd(' ') # Try to reset to terminal default for day


    def _spawn_obstacle(self):
        # CORRECTED: Use OBSTACLE_SPRITES
        sprite_lines = random.choice(OBSTACLE_SPRITES)

        # Ensure obstacles spawn off-screen to the right
        # x_pos is for the left edge of the sprite
        obstacle_width = max(len(line) for line in sprite_lines) if sprite_lines else 0
        min_x = self.width + int(self.current_obstacle_speed * 2)
        max_x = self.width + int(self.current_obstacle_speed * 5) + obstacle_width
        x_pos = random.randint(min_x, max_x)

        # CORRECTED: Pass sprite_lines and use self.ground_line_y
        self.obstacles.append(Obstacle(x_pos, self.ground_line_y, sprite_lines, self.current_obstacle_speed))


    def _handle_input(self, key):
        if self.game_state == "PLAYING":
            if key == ord(' ') or key == curses.KEY_UP:
                self.player.jump()
            elif key == ord('p') or key == ord('P'): # Accept lowercase or uppercase
                self.game_state = "PAUSED"
            elif key == ord('q') or key == ord('Q'):
                return False # Signal to quit
        elif self.game_state == "PAUSED":
            if key == ord('p') or key == ord('P'):
                self.game_state = "PLAYING"
            elif key == ord('q') or key == ord('Q'):
                return False
        elif self.game_state == "GAME_OVER":
            if key == ord('r') or key == ord('R'):
                self._reset_game()
            elif key == ord('q') or key == ord('Q'):
                return False
        elif self.game_state == "START":
            if key != -1: # Any key to start
                self._reset_game()
        return True

    def _update_game_state(self):
        if self.game_state != "PLAYING":
            return

        self.frame_count += 1
        self.score = self.frame_count // (FPS // 2) # Score increases a bit faster

        if self.score > self.high_score:
            self.high_score = self.score

        self.stage_timer_frames += 1
        self.day_night_timer_frames += 1

        if self.day_night_timer_frames >= DAY_NIGHT_CYCLE_SEC * FPS:
            self.is_day_mode = not self.is_day_mode
            self.day_night_timer_frames = 0
            bg_color_to_set = curses.color_pair(COLOR_PAIR_NIGHT_GROUND_TEXT) if not self.is_day_mode else curses.color_pair(COLOR_PAIR_DAY_GROUND_TEXT)
            try:
                self.stdscr.bkgd(' ', bg_color_to_set)
            except curses.error:
                if self.is_day_mode: self.stdscr.bkgd(' ')


        if self.stage_timer_frames >= STAGE_DURATION_SEC * FPS:
            self.stage_timer_frames = 0
            self.current_obstacle_speed += SPEED_INCREASE_PER_STAGE
            new_spawn_interval_sec = (self.current_obstacle_spawn_interval_frames / FPS) - SPAWN_DECREASE_PER_STAGE_SEC
            new_spawn_interval_sec = max(MIN_OBSTACLE_SPAWN_INTERVAL_SEC, new_spawn_interval_sec)
            self.current_obstacle_spawn_interval_frames = int(new_spawn_interval_sec * FPS)

        self.player.update()

        if self.frame_count - self.last_obstacle_spawn_frame >= self.current_obstacle_spawn_interval_frames:
            can_spawn = True
            if self.obstacles:
                last_obs = self.obstacles[-1]
                # A simple threshold: don't spawn if the last obstacle hasn't moved far enough from the right edge
                spawn_clearance = last_obs.width * 2.5 # Increased clearance a bit
                if last_obs.x + last_obs.width > self.width - spawn_clearance:
                    can_spawn = False
            if can_spawn:
                self._spawn_obstacle()
                self.last_obstacle_spawn_frame = self.frame_count

        new_obstacles = []
        for obs in self.obstacles:
            obs.update()
            if not obs.is_offscreen():
                new_obstacles.append(obs)
        self.obstacles = new_obstacles

        player_rect = self.player.get_rect()
        for obs in self.obstacles:
            obs_rect = obs.get_rect()
            # AABB collision
            if (player_rect[0] < obs_rect[0] + obs_rect[2] and
                    player_rect[0] + player_rect[2] > obs_rect[0] and
                    player_rect[1] < obs_rect[1] + obs_rect[3] and
                    player_rect[1] + player_rect[3] > obs_rect[1]):
                self.game_state = "GAME_OVER"
                self._save_high_score()
                break

    def _render(self):
        self.stdscr.erase()
        current_bg_color = self._get_color_pair("ground_text")

        # Explicitly set background for the entire screen based on mode
        # This helps maintain consistency, especially for day mode with default backgrounds
        try:
            self.stdscr.bkgd(' ', current_bg_color)
        except curses.error:
            if self.is_day_mode: # Fallback for day mode if specific bkgd fails
                self.stdscr.bkgd(' ')


        # Draw Ground Line (visual line is one below where feet touch)
        # The actual alignment point is self.ground_line_y
        for x_coord in range(self.width -1 ): # -1 to avoid curses error at far right edge
            try:
                self.stdscr.addch(self.ground_line_y + 1, x_coord, GROUND_CHAR, current_bg_color)
            except curses.error: pass

        # CORRECTED: Draw Player Sprite
        player_sprite_lines = self.player.get_current_sprite_lines()
        px, py_top, pw, ph = self.player.get_rect() # player_x, player_y_top, player_width, player_height
        dino_color = self._get_color_pair("dino")
        for i, line_str in enumerate(player_sprite_lines):
            draw_y = py_top + i
            if 0 <= draw_y < self.height: # Check if row is on screen
                for j, char_to_draw in enumerate(line_str):
                    if char_to_draw != ' ': # Simple transparency: don't draw spaces
                        draw_x = px + j
                        if 0 <= draw_x < self.width: # Check if col is on screen
                            try:
                                self.stdscr.addch(draw_y, draw_x, char_to_draw, dino_color)
                            except curses.error: pass # Ignore errors drawing at edges

        # CORRECTED: Draw Obstacle Sprites
        obstacle_color = self._get_color_pair("obstacle")
        for obs in self.obstacles:
            ox, oy_top, ow, oh = obs.get_rect() # obs_x_left, obs_y_top, obs_width, obs_height
            for i, line_str in enumerate(obs.sprite_lines):
                draw_y = oy_top + i
                if 0 <= draw_y < self.height:
                    for j, char_to_draw in enumerate(line_str):
                        if char_to_draw != ' ': # Simple transparency
                            # Round obs.x for drawing as it can be a float due to speed
                            draw_x = int(round(ox)) + j
                            if 0 <= draw_x < self.width:
                                try:
                                    self.stdscr.addch(draw_y, draw_x, char_to_draw, obstacle_color)
                                except curses.error: pass

        # Draw Score and Info
        score_text = f"Score: {self.score}"
        hs_text = f"High Score: {self.high_score}"
        mode_text = "Mode: Day" if self.is_day_mode else "Mode: Night"
        try:
            self.stdscr.addstr(0, 1, score_text, current_bg_color)
            # Adjust HS text position if too close to right edge
            hs_text_x = self.width - len(hs_text) - 1
            self.stdscr.addstr(0, hs_text_x if hs_text_x > len(score_text) + 2 else (self.width - len(hs_text)) // 2, hs_text, current_bg_color)
            self.stdscr.addstr(1, 1, mode_text, current_bg_color)
        except curses.error: pass


        if self.game_state == "START":
            msg_lines = ["Command-Line Dino Runner", " ", "Press any key to start", "Space/Up: Jump, P: Pause, Q: Quit"]
            start_y = self.height // 2 - len(msg_lines) // 2
            for i, line in enumerate(msg_lines):
                try: self.stdscr.addstr(start_y + i, max(0,(self.width - len(line)) // 2), line, current_bg_color)
                except: pass # Ignore errors if text doesn't fit
        elif self.game_state == "GAME_OVER":
            msg_lines = ["GAME OVER", f"Your score: {self.score}", "Press r to restart, q to quit"]
            start_y = self.height // 2 - len(msg_lines) // 2
            for i, line in enumerate(msg_lines):
                try: self.stdscr.addstr(start_y + i, max(0,(self.width - len(line)) // 2), line, current_bg_color)
                except: pass
        elif self.game_state == "PAUSED":
            msg = "-- PAUSED --"
            try: self.stdscr.addstr(self.height // 2, max(0,(self.width - len(msg)) // 2), msg, current_bg_color)
            except: pass

        self.stdscr.refresh()

    def run(self):
        curses.curs_set(0)
        self.stdscr.nodelay(True)
        self.stdscr.keypad(True)

        running = True
        while running:
            frame_start_time = time.monotonic()
            key = self.stdscr.getch() # Get input
            running = self._handle_input(key) # Process input
            if not running: break

            self._update_game_state() # Update game logic
            self._render() # Draw the current frame

            # Frame rate control
            elapsed_time = time.monotonic() - frame_start_time
            sleep_time = TARGET_FRAME_TIME - elapsed_time
            if sleep_time > 0:
                curses.napms(int(sleep_time * 1000))
        self._save_high_score()

def main_game_loop(stdscr):
    # Check for color support
    if not curses.has_colors():
        # This is a basic fallback. Ideally, you'd have a truly monochrome mode
        # or inform the user more gracefully. For now, it will try to run.
        pass # init_pair calls might silently fail or use defaults.

    # Check terminal size
    min_height, min_width = 20, 60 # Minimum dimensions for sprites to look okay
    h, w = stdscr.getmaxyx()
    if h < min_height or w < min_width:
        stdscr.clear()
        try:
            msg = f"Terminal too small ({w}x{h}). Need {min_width}x{min_height}."
            stdscr.addstr(0,0, msg)
            stdscr.addstr(1,0, "Please resize and restart.")
        except curses.error: # If even this fails (extremely small terminal)
            pass # Just exit gracefully
        stdscr.refresh()
        stdscr.getch() # Wait for a key press before exiting
        return

    # Initial background setting based on potential day mode start
    # This ensures the screen is cleared with the correct initial colors
    # before the first frame if there's any delay.
    if curses.has_colors():
        curses.start_color() # Ensure colors are started
        curses.use_default_colors()
        try:
            # Attempt to set initial background for day mode (white on default)
            # This helps clear the screen correctly from the start.
            curses.init_pair(COLOR_PAIR_DAY_GROUND_TEXT, curses.COLOR_WHITE, -1)
            stdscr.bkgd(' ', curses.color_pair(COLOR_PAIR_DAY_GROUND_TEXT))
        except curses.error:
            stdscr.bkgd(' ') # Fallback to simple clear if bkgd with pair fails

    game = Game(stdscr)
    game.run()


if __name__ == "__main__":
    try:
        curses.wrapper(main_game_loop)
    except curses.error as e:
        print(f"Curses error: {e}")
        print("Your terminal might be too small, not support colors, or lack other capabilities.")
        print("Try resizing your terminal or using a different one (e.g., Windows Terminal, WSL with a modern terminal).")
    except Exception as e:
        import traceback
        print(f"An unexpected error occurred: {e}")
        print(traceback.format_exc()) # Print full traceback for debugging
    finally:
        print("Game exited. Thanks for playing!")
