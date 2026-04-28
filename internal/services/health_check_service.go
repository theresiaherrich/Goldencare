package service

import (
	"fmt"
	"strings"

	"github.com/theresiaherrich/Goldencare/internal/models"
)

type HealthStatus string

const (
	StatusStabil    HealthStatus = "STABIL"
	StatusObservasi HealthStatus = "OBSERVASI"
	StatusWaspada   HealthStatus = "WASPADA"
	StatusDarurat   HealthStatus = "DARURAT"
)

type HealthPhaseColor string

const (
	ColorGreen  HealthPhaseColor = "hijau"
	ColorYellow HealthPhaseColor = "kuning"
	ColorOrange HealthPhaseColor = "oranye"
	ColorRed    HealthPhaseColor = "merah"
)

type HealthRecommendation struct {
	Status           HealthStatus     `json:"status"`
	Color            HealthPhaseColor `json:"color"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	ActionItems      []string         `json:"action_items"`
	AnomaliesFound   []string         `json:"anomalies_found"`
	CriticalFactors  []string         `json:"critical_factors,omitempty"`
	FollowUpInterval string           `json:"follow_up_interval,omitempty"`
}

type VitalSign struct {
	SystolicBP  int
	DiastolicBP int
	HeartRate   int
	Temperature float64
	BloodSugar  float64
	Weight      float64
}

func AnalyzeHealthCondition(req models.CreatePemeriksaanRequest) HealthRecommendation {
	tandaVital := parseVitalSigns(req)
	tingkatKeparahan := determineOverallSeverity(tandaVital)

	switch tingkatKeparahan {
	case StatusStabil:
		return buildStabilRecommendation(tandaVital)
	case StatusObservasi:
		return buildObservasiRecommendation(tandaVital)
	case StatusWaspada:
		return buildWaspadaRecommendation(tandaVital)
	case StatusDarurat:
		return buildDaruratRecommendation(tandaVital)
	default:
		return buildStabilRecommendation(tandaVital)
	}
}

func parseVitalSigns(req models.CreatePemeriksaanRequest) VitalSign {
	tandaVital := VitalSign{
		HeartRate:   req.DetakJantung,
		Temperature: req.SuhuTubuh,
		BloodSugar:  req.GulaDarah,
		Weight:      req.BeratBadan,
	}

	if req.TekananDarah != "" {
		bagian := strings.Split(req.TekananDarah, "/")
		if len(bagian) == 2 {
			fmt.Sscanf(strings.TrimSpace(bagian[0]), "%d", &tandaVital.SystolicBP)
			fmt.Sscanf(strings.TrimSpace(bagian[1]), "%d", &tandaVital.DiastolicBP)
		}
	}

	return tandaVital
}

func determineOverallSeverity(tandaVital VitalSign) HealthStatus {
	if isDarurat(tandaVital) {
		return StatusDarurat
	}
	if isWaspada(tandaVital) {
		return StatusWaspada
	}
	if isObservasi(tandaVital) {
		return StatusObservasi
	}
	return StatusStabil
}

func isDarurat(tandaVital VitalSign) bool {
	if tandaVital.SystolicBP > 0 && tandaVital.SystolicBP >= 180 {
		return true
	}
	if tandaVital.DiastolicBP > 0 && tandaVital.DiastolicBP >= 120 {
		return true
	}

	if tandaVital.HeartRate > 0 && (tandaVital.HeartRate < 40 || tandaVital.HeartRate > 130) {
		return true
	}

	if tandaVital.Temperature > 0 && tandaVital.Temperature >= 39.0 {
		return true
	}

	if tandaVital.BloodSugar > 0 && (tandaVital.BloodSugar > 250 || tandaVital.BloodSugar < 55) {
		return true
	}

	return false
}

func isWaspada(tandaVital VitalSign) bool {
	if tandaVital.SystolicBP > 0 && tandaVital.SystolicBP >= 151 && tandaVital.SystolicBP <= 179 {
		return true
	}
	if tandaVital.DiastolicBP > 0 && tandaVital.DiastolicBP >= 91 && tandaVital.DiastolicBP <= 119 {
		return true
	}

	if tandaVital.HeartRate > 0 && ((tandaVital.HeartRate >= 40 && tandaVital.HeartRate <= 49) || (tandaVital.HeartRate >= 111 && tandaVital.HeartRate <= 130)) {
		return true
	}

	if tandaVital.Temperature > 0 && ((tandaVital.Temperature >= 37.9 && tandaVital.Temperature <= 38.9) || tandaVital.Temperature < 35.5) {
		return true
	}

	if tandaVital.BloodSugar > 0 && ((tandaVital.BloodSugar > 180 && tandaVital.BloodSugar <= 250) || (tandaVital.BloodSugar >= 55 && tandaVital.BloodSugar <= 69)) {
		return true
	}

	return false
}

func isObservasi(tandaVital VitalSign) bool {
	if tandaVital.SystolicBP > 0 && ((tandaVital.SystolicBP >= 131 && tandaVital.SystolicBP <= 150) || (tandaVital.SystolicBP >= 80 && tandaVital.SystolicBP <= 89)) {
		return true
	}

	if tandaVital.HeartRate > 0 && ((tandaVital.HeartRate >= 50 && tandaVital.HeartRate <= 59) || (tandaVital.HeartRate >= 101 && tandaVital.HeartRate <= 110)) {
		return true
	}

	if tandaVital.Temperature > 0 && tandaVital.Temperature >= 37.3 && tandaVital.Temperature <= 37.8 {
		return true
	}

	if tandaVital.BloodSugar > 0 && ((tandaVital.BloodSugar > 140 && tandaVital.BloodSugar <= 180) || (tandaVital.BloodSugar >= 70 && tandaVital.BloodSugar <= 79)) {
		return true
	}

	return false
}

func buildStabilRecommendation(tandaVital VitalSign) HealthRecommendation {
	return HealthRecommendation{
		Status:      StatusStabil,
		Color:       ColorGreen,
		Title:       "Kondisi Vital Normal",
		Description: "Tanda vital lansia berada dalam batas aman. Lanjutkan asuhan perawatan harian sesuai jadwal. Pastikan asupan gizi dan cairan harian tetap terpenuhi.",
		ActionItems: []string{
			"Lanjutkan asuhan perawatan harian sesuai jadwal",
			"Pastikan asupan gizi dan cairan harian tetap terpenuhi",
			"Pertahankan rutinitas aktivitas harian",
			"Lakukan pemeriksaan berkala sesuai jadwal",
		},
		AnomaliesFound: []string{},
	}
}

func buildObservasiRecommendation(tandaVital VitalSign) HealthRecommendation {
	anomali := findAnomalies(tandaVital, StatusObservasi)

	jenisBencana := "Kondisi"
	for _, a := range anomali {
		if strings.Contains(a, "Suhu") {
			jenisBencana = "Suhu Tubuh"
			break
		} else if strings.Contains(a, "Tekanan") {
			jenisBencana = "Tensi"
			break
		} else if strings.Contains(a, "Detak") {
			jenisBencana = "Nadi"
			break
		} else if strings.Contains(a, "Gula") {
			jenisBencana = "Gula Darah"
			break
		}
	}

	return HealthRecommendation{
		Status:      StatusObservasi,
		Color:       ColorYellow,
		Title:       fmt.Sprintf("Anomali Ringan pada %s", jenisBencana),
		Description: "Terdapat sedikit anomali pada hasil pemeriksaan. Lakukan pemantauan berkala setiap 2–4 jam ke depan. Anjurkan lansia untuk lebih banyak beristirahat dan minum air hangat.",
		ActionItems: []string{
			"Lakukan pemantauan berkala setiap 2-4 jam ke depan",
			"Anjurkan lansia untuk lebih banyak beristirahat",
			"Pastikan lansia minum air hangat secara teratur",
			"Catat perkembangan kondisi setiap 2 jam",
			"Hubungi perawat jika ada perubahan signifikan",
		},
		AnomaliesFound:   anomali,
		FollowUpInterval: "2-4 jam",
	}
}

func buildWaspadaRecommendation(tandaVital VitalSign) HealthRecommendation {
	anomali := findAnomalies(tandaVital, StatusWaspada)

	return HealthRecommendation{
		Status:      StatusWaspada,
		Color:       ColorOrange,
		Title:       "Perlu Intervensi Medis",
		Description: "Peringatan: Indikator kesehatan berada di luar batas aman. Berikan obat pereda/penurun sesuai resep dokter jika ada, dan segera laporkan kondisi ini kepada Perawat Kepala panti.",
		ActionItems: []string{
			"Berikan obat sesuai resep dokter jika tersedia",
			"Segera laporkan kondisi ke Perawat Kepala panti",
			"Lakukan pemantauan intensif setiap 1 jam",
			"Siapkan data vital untuk konsultasi dokter",
			"Dokumentasikan semua perubahan kondisi",
			"Pastikan lansia istirahat yang cukup",
		},
		AnomaliesFound:   anomali,
		CriticalFactors:  anomali,
		FollowUpInterval: "1 jam",
	}
}

func buildDaruratRecommendation(tandaVital VitalSign) HealthRecommendation {
	anomali := findAnomalies(tandaVital, StatusDarurat)

	return HealthRecommendation{
		Status:      StatusDarurat,
		Color:       ColorRed,
		Title:       "KONDISI KRITIS!",
		Description: "Segera lakukan prosedur Pertolongan Pertama (First Aid)! Kondisi ini mengancam jiwa. Segera hubungi dokter panti atau siapkan rujukan ambulans ke Rumah Sakit terdekat.",
		ActionItems: []string{
			"⚠️ SEGERA HUBUNGI DOKTER/AMBULANS",
			"Lakukan First Aid sesuai prosedur panti",
			"Siapkan rujukan ke Rumah Sakit terdekat",
			"Dokumentasikan waktu dan kondisi kritis",
			"Hubungi keluarga lansia segera",
			"Jangan tinggalkan lansia sendirian",
		},
		AnomaliesFound:   anomali,
		CriticalFactors:  anomali,
		FollowUpInterval: "IMMEDIATE",
	}
}

func findAnomalies(tandaVital VitalSign, tingkatKeparahan HealthStatus) []string {
	var anomali []string

	if tandaVital.SystolicBP > 0 || tandaVital.DiastolicBP > 0 {
		switch tingkatKeparahan {
		case StatusDarurat:
			if tandaVital.SystolicBP >= 180 || tandaVital.DiastolicBP >= 120 {
				anomali = append(anomali, fmt.Sprintf("Tekanan Darah: Krisis Hipertensi (%d/%d mmHg)", tandaVital.SystolicBP, tandaVital.DiastolicBP))
			}
		case StatusWaspada:
			if (tandaVital.SystolicBP >= 151 && tandaVital.SystolicBP <= 179) || (tandaVital.DiastolicBP >= 91 && tandaVital.DiastolicBP <= 119) {
				anomali = append(anomali, fmt.Sprintf("Tekanan Darah: Hipertensi Tahap 2 (%d/%d mmHg)", tandaVital.SystolicBP, tandaVital.DiastolicBP))
			}
		case StatusObservasi:
			if tandaVital.SystolicBP >= 131 && tandaVital.SystolicBP <= 150 {
				anomali = append(anomali, fmt.Sprintf("Tekanan Darah: Pra-hipertensi Sistolik %d mmHg", tandaVital.SystolicBP))
			} else if tandaVital.SystolicBP >= 80 && tandaVital.SystolicBP <= 89 {
				anomali = append(anomali, fmt.Sprintf("Tekanan Darah: Hipotensi ringan Sistolik %d mmHg", tandaVital.SystolicBP))
			}
		}
	}

	if tandaVital.HeartRate > 0 {
		switch tingkatKeparahan {
		case StatusDarurat:
			if tandaVital.HeartRate < 40 || tandaVital.HeartRate > 130 {
				anomali = append(anomali, fmt.Sprintf("Detak Jantung: Kritis (%d BPM)", tandaVital.HeartRate))
			}
		case StatusWaspada:
			if (tandaVital.HeartRate >= 40 && tandaVital.HeartRate <= 49) || (tandaVital.HeartRate >= 111 && tandaVital.HeartRate <= 130) {
				anomali = append(anomali, fmt.Sprintf("Detak Jantung: Tidak Normal (%d BPM)", tandaVital.HeartRate))
			}
		case StatusObservasi:
			if tandaVital.HeartRate >= 50 && tandaVital.HeartRate <= 59 {
				anomali = append(anomali, fmt.Sprintf("Detak Jantung: Agak lambat (%d BPM)", tandaVital.HeartRate))
			} else if tandaVital.HeartRate >= 101 && tandaVital.HeartRate <= 110 {
				anomali = append(anomali, fmt.Sprintf("Detak Jantung: Agak cepat (%d BPM)", tandaVital.HeartRate))
			}
		}
	}

	if tandaVital.Temperature > 0 {
		switch tingkatKeparahan {
		case StatusDarurat:
			if tandaVital.Temperature >= 39.0 {
				anomali = append(anomali, fmt.Sprintf("Suhu Tubuh: Demam Tinggi (%.1f°C)", tandaVital.Temperature))
			}
		case StatusWaspada:
			if (tandaVital.Temperature >= 37.9 && tandaVital.Temperature <= 38.9) || tandaVital.Temperature < 35.5 {
				anomali = append(anomali, fmt.Sprintf("Suhu Tubuh: Abnormal (%.1f°C)", tandaVital.Temperature))
			}
		case StatusObservasi:
			if tandaVital.Temperature >= 37.3 && tandaVital.Temperature <= 37.8 {
				anomali = append(anomali, fmt.Sprintf("Suhu Tubuh: Subfebris (%.1f°C)", tandaVital.Temperature))
			}
		}
	}

	if tandaVital.BloodSugar > 0 {
		switch tingkatKeparahan {
		case StatusDarurat:
			if tandaVital.BloodSugar > 250 || tandaVital.BloodSugar < 55 {
				anomali = append(anomali, fmt.Sprintf("Gula Darah: Kritis (%.1f mg/dL)", tandaVital.BloodSugar))
			}
		case StatusWaspada:
			if (tandaVital.BloodSugar > 180 && tandaVital.BloodSugar <= 250) || (tandaVital.BloodSugar >= 55 && tandaVital.BloodSugar <= 69) {
				anomali = append(anomali, fmt.Sprintf("Gula Darah: Tidak Normal (%.1f mg/dL)", tandaVital.BloodSugar))
			}
		case StatusObservasi:
			if (tandaVital.BloodSugar > 140 && tandaVital.BloodSugar <= 180) || (tandaVital.BloodSugar >= 70 && tandaVital.BloodSugar <= 79) {
				anomali = append(anomali, fmt.Sprintf("Gula Darah: Agak Tinggi/Rendah (%.1f mg/dL)", tandaVital.BloodSugar))
			}
		}
	}

	return anomali
}
