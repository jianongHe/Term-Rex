package game

import "time"

// applyStage smoothly transitions parameters based on score threshold crossings.
func (g *Game) applyStage() {
	// determine target stage for current score
	target := 0
	for i := len(stageConfigs) - 1; i >= 0; i-- {
		if g.score >= stageConfigs[i].ScoreThreshold {
			target = i
			break
		}
	}
	// on first crossing, start transition
	if target != g.stageIndexTarget {
		g.stageIndexTarget = target
		g.stageTransitionStart = time.Now()

		// 当进入新阶段时，触发分数闪烁效果
		if target > g.stageIndexActive {
			g.scoreBlinking = true
			g.scoreBlinkStart = time.Now()
			g.scoreBlinkVisible = true
			g.lastBlinkToggle = time.Now()

			// 播放得分音效
			GetAudioManager().PlaySound(SoundScore)
		}
	}
	// if currently transitioning between two stages
	if g.stageIndexActive != g.stageIndexTarget {
		elapsed := time.Since(g.stageTransitionStart)
		frac := float64(elapsed) / float64(stageTransitionDuration)
		if frac >= 1 {
			// finish transition
			g.stageIndexActive = g.stageIndexTarget
			obstacleSpeed = stageConfigs[g.stageIndexActive].Speed * speedFactor
			birdProbability = stageConfigs[g.stageIndexActive].BirdProb
			bigBirdProbability = stageConfigs[g.stageIndexActive].BigBirdProb
			groupCactusProbability = stageConfigs[g.stageIndexActive].GroupCactusProb
			g.stageTransitionStart = time.Time{}
		} else {
			// interpolate between active and target
			old := stageConfigs[g.stageIndexActive]
			next := stageConfigs[g.stageIndexTarget]
			speed := old.Speed + frac*(next.Speed-old.Speed)
			obstacleSpeed = speed * speedFactor
			birdProbability = old.BirdProb + frac*(next.BirdProb-old.BirdProb)
			bigBirdProbability = old.BigBirdProb + frac*(next.BigBirdProb-old.BigBirdProb)
			groupCactusProbability = old.GroupCactusProb + frac*(next.GroupCactusProb-old.GroupCactusProb)
		}
	} else {
		// no transition: keep active stage values
		sc := stageConfigs[g.stageIndexActive]
		obstacleSpeed = sc.Speed * speedFactor
		birdProbability = sc.BirdProb
		bigBirdProbability = sc.BigBirdProb
		groupCactusProbability = sc.GroupCactusProb
	}
}
