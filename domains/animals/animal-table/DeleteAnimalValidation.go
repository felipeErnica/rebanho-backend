package animalTable

import (
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type WarnRecordValidations struct {
	HasLactation          bool `db:"has_lactation"`
	HasSlaughter          bool `db:"has_slaughter"`
	HasInsemination       bool `db:"has_insemination"`
	HasBreeding           bool `db:"has_breeding"`
	HasTransfer           bool `db:"has_transfer"`
	HasBullPastureEntries bool `db:"has_bull_pastures_entries"`
}

type ErrorRecordValidations struct {
	HasChildren        bool `db:"has_children"`
	IsCalfInLactation  bool `db:"is_calf_lac"`
	IsEmbryoDonor      bool `db:"is_embryo_donor"`
	IsEmbryoBull       bool `db:"is_embryo_bull"`
	IsInseminationBull bool `db:"is_insemination_bull"`
	IsBreedingBull     bool `db:"is_breeding_bull"`
}

const DELETE_OBSERVATION = "\n\nOBS.: A exclusão de animais só é recomendada em caso de erros." +
	"Em caso de morte e/ou abate, faça o registro apropriado." +
	"Lembre-se, apagar um animal do sistema apagará, PERMANENTEMENTE, todas as informções ligadas a ele."

func validDelete(db *sqlx.DB, id string, userId string) *apiError.APIError {

	err := throwError(db, id, userId)
	if err != nil {
		return err
	}

	err = hasImportantRecords(db, id, userId)
	if err != nil {
		return err
	}

	return nil
}

func throwError(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select 
			exists (
				select 1
				from animals
				where deleted_at is null
					and (mother_id = $1 or father_id = $1)
					and user_id = $2
			) as has_children,
			exists (
				select 1
				from lactations
				where deleted_at is null
					and calf_id = $1
					and user_id = $2
			) as is_calf_lac,
			exists (
				select 1
				from embryo_transfer
				where donor_id = $1
					and user_id = $2
					and deleted_at is null
			) as is_embryo_donor,
			exists (
				select 1
				from embryo_transfer
				where bull_id = $1
					and user_id = $2
					and deleted_at is null
			) as is_embryo_bull,
			exists (
				select 1
				from insemination_entries
				where bull_id = $1
					and user_id = $2
					and deleted_at is null
			) as is_insemination_bull,
			exists (
				select 1
				from natural_breedings
				where bull_id = $1
					and user_id = $2
					and deleted_at is null
			) as is_breeding_bull
	`

	res, err := repositoriesUtil.GetOne[ErrorRecordValidations](db, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	messages := []string{}
	if res.HasChildren {
		messages = append(messages, "O animal possui filhos registrados no sistema.")
	}

	if res.IsCalfInLactation {
		messages = append(messages, "O animal está ligado, como bezerro, a uma lactação.")
	}

	if res.IsEmbryoDonor {
		messages = append(messages, "A vaca é uma doadora de embriões.")
	}

	if res.IsInseminationBull {
		messages = append(messages, "O touro possui registros ativos de inseminação.")
	}

	if res.IsBreedingBull {
		messages = append(messages, "O touro possui registros ativos de cobertura.")
	}

	if res.IsEmbryoBull {
		messages = append(messages, "O touro possui registros ativos de transferência embrionária.")
	}

	if len(messages) != 0 {
		warnMsg := "O animal não pode ser excluído, devido aos seguintes motivos:"
		formatedMsg := []string{}
		for i, message := range messages {
			msg := fmt.Sprintf("%d - %s", i+1, message)
			formatedMsg = append(formatedMsg, msg)
		}
		resultMsg := strings.Join(formatedMsg, "\n")
		return apiError.DeleteAPIError(warnMsg + "\n" + resultMsg + DELETE_OBSERVATION)
	}

	return nil
}

func hasImportantRecords(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select 
			exists (
				select 1
				from lactations
				where deleted_at is null
					and animal_id = $1
					and user_id = $2
			) as has_lactation,
			exists (
				select 1
				from natural_breedings
				where deleted_at is null
					and animal_id = $1
					and user_id = $2
			) as has_breeding,
			exists (
				select 1
				from embryo_transfer
				where deleted_at is null
					and receiver_id = $1
					and user_id = $2
			) as has_transfer,
			exists (
				select 1
				from slaughter_entries
				where deleted_at is null
					and animal_id = $1
					and user_id = $2
			) as has_slaughter,
			exists (
				select 1
				from insemination_entries
				where deleted_at is null
					and animal_id = $1
					and user_id = $2
			) as has_insemination,
			exists (
				select 1
				from pasture_entries pe
					join animals a on a.id = pe.animal_id
				where pe.deleted_at is null
					and a.sex = 'M'
					and a.animal_type = 'REPRODUCTION_ANIMAL'
					and pe.animal_id = $1
					and pe.user_id = $2
			) as has_bull_pastures_entries
	`

	res, err := repositoriesUtil.GetOne[WarnRecordValidations](db, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	records := make([]string, 0)
	if res.HasLactation {
		records = append(records, "lactações")
	}

	if res.HasBreeding {
		records = append(records, "coberturas")
	}

	if res.HasInsemination {
		records = append(records, "inseminações")
	}

	if res.HasSlaughter {
		records = append(records, "abate")
	}

	if res.HasTransfer {
		records = append(records, "receptação de embrião")
	}

	if res.HasBullPastureEntries {
		records = append(records, "entradas em lotes")
	}

	if len(records) != 0 {
		recordStr := strings.Join(records, ", ")
		warnMsg := fmt.Sprintf("Este animal possui importantes registros de: %s. Deseja exclui-lo mesmo assim?", recordStr)
		return apiError.DeleteWarning(warnMsg + DELETE_OBSERVATION)
	}

	return nil
}
