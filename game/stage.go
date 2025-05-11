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

			// 设置当前阶段的速度
			obstacleSpeed = stageConfigs[g.stageIndexActive].Speed * speedFactor

			// 更新障碍物间距和概率
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

			// 平滑过渡速度
			speed := old.Speed + frac*(next.Speed-old.Speed)
			obstacleSpeed = speed * speedFactor

			// 平滑过渡障碍物间距
			minGap := int(float64(old.MinGap) + frac*float64(next.MinGap-old.MinGap))
			maxGap := int(float64(old.MaxGap) + frac*float64(next.MaxGap-old.MaxGap))

			// 平滑过渡概率
			cactusProbability := old.CactusProb + frac*(next.CactusProb-old.CactusProb)
			singleCactusRatio := old.SingleCactusRatio + frac*(next.SingleCactusRatio-old.SingleCactusRatio)
			shortCactusRatio := old.ShortCactusRatio + frac*(next.ShortCactusRatio-old.ShortCactusRatio)
			groupCactusRatio := old.GroupCactusRatio + frac*(next.GroupCactusRatio-old.GroupCactusRatio)
			smallBirdRatio := old.SmallBirdRatio + frac*(next.SmallBirdRatio-old.SmallBirdRatio)
			bigBirdRatio := old.BigBirdRatio + frac*(next.BigBirdRatio-old.BigBirdRatio)

			// 创建临时阶段配置
			tempStage := StageConfig{
				CactusProb:        cactusProbability,
				SingleCactusRatio: singleCactusRatio,
				ShortCactusRatio:  shortCactusRatio,
				GroupCactusRatio:  groupCactusRatio,
				SmallBirdRatio:    smallBirdRatio,
				BigBirdRatio:      bigBirdRatio,
			}

			// 更新障碍物管理器的概率
			g.obstacleManager.cactusProbability = tempStage.CactusProb
			g.obstacleManager.singleCactusRatio = tempStage.SingleCactusRatio
			g.obstacleManager.shortCactusRatio = tempStage.ShortCactusRatio
			g.obstacleManager.groupCactusRatio = tempStage.GroupCactusRatio
			g.obstacleManager.smallBirdRatio = tempStage.SmallBirdRatio
			g.obstacleManager.bigBirdRatio = tempStage.BigBirdRatio

			// 更新障碍物间距
			g.obstacleManager.minGap = minGap
			g.obstacleManager.maxGap = maxGap
		}
	} else {
		// no transition: keep active stage values
		sc := stageConfigs[g.stageIndexActive]
		obstacleSpeed = sc.Speed * speedFactor

		// 确保障碍物间距与当前阶段一致
		g.obstacleManager.UpdateStageGaps(sc.MinGap, sc.MaxGap, g.stageIndexActive)
	}
}
