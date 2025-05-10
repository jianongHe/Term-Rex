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

			// 更新障碍物间距
			g.obstacleManager.UpdateStageGaps(
				stageConfigs[g.stageIndexActive].MinGap,
				stageConfigs[g.stageIndexActive].MaxGap,
				g.stageIndexActive,
			)

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

			// 平滑过渡障碍物间距
			minGap := int(float64(old.MinGap) + frac*float64(next.MinGap-old.MinGap))
			maxGap := int(float64(old.MaxGap) + frac*float64(next.MaxGap-old.MaxGap))
			stageIndex := g.stageIndexActive // 在过渡期间使用当前活动阶段
			g.obstacleManager.UpdateStageGaps(minGap, maxGap, stageIndex)
		}
	} else {
		// no transition: keep active stage values
		sc := stageConfigs[g.stageIndexActive]
		obstacleSpeed = sc.Speed * speedFactor
		birdProbability = sc.BirdProb
		bigBirdProbability = sc.BigBirdProb
		groupCactusProbability = sc.GroupCactusProb

		// 确保障碍物间距与当前阶段一致
		g.obstacleManager.UpdateStageGaps(sc.MinGap, sc.MaxGap, g.stageIndexActive)
	}
}
